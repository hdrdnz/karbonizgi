package controllers

import (
	"carbonfootprint/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Suggest struct {
	Action string `json:"action"`
	Status string `json:"status"`
}

type SuggestIds struct {
	ActionId int    `json:"action_id"`
	Action   string `json:"action"`
	Status   string `json:"status"`
}

type UserSug struct {
	Status string       `json:"status"`
	Data   []SuggestIds `json:"data"`
}

// @Summary      Kullanıcı aksiyon
// @Description  Kullanıcı Aksiyon Kayıtları
// @Tags         Action
// @Accept       json
// @Produce      json
// @Success      200 {object} UserSug "success"
// @Router       /user-suggest [get]
func UserSuggest(c *gin.Context) {
	db := model.GetDB()
	user := model.User{}

	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı bulunamadı.",
		})
		return
	}
	actions := []model.UserAction{}
	if err := db.Where("user_id=?", user.Id).Find(&actions).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}

	if len(actions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcıya ait aksiyon bulunmamaktadır.",
		})
		return
	}
	var userAct []SuggestIds
	for _, act := range actions {
		userAct = append(userAct, SuggestIds{
			ActionId: act.Id,
			Action:   act.Action,
			Status:   act.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   userAct,
	})

}

// @Summary      Kullanıcı aksiyon ekleme
// @Description  Aksiyon Ekleme Kısmı
// @Tags         Action
// @Accept       json
// @Produce      json
// @Param        Suggest body Suggest true "status değeri için sadece 'planned,in_progress,completed,cancelled' değerlerden biri gönderilmelidir. "
// @Success      200 {object} Response "success"
// @Failure      400 {object} Response "Invalid request"
// @Router       /add-suggest [post]
func AddSuggest(c *gin.Context) {
	db := model.GetDB()
	user := model.User{}
	var suggest Suggest
	if err := c.ShouldBindJSON(&suggest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı bulunamadı.",
		})
		return
	}

	if suggest.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Aksiyon boş olamaz.",
		})
		return
	}
	if suggest.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Status boş olamaz.",
		})
		return
	}
	isStatus := false
	status := []string{"planned", "in_progress", "completed", "cancelled"}
	for _, st := range status {
		if suggest.Status == st {
			isStatus = true
			break
		}
	}
	if !isStatus {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz aksiyon status ifadesi girilmiştir.",
		})
		return
	}

	userAct := model.UserAction{}
	userAct.UserId = user.Id
	userAct.Action = suggest.Action
	userAct.Status = suggest.Status
	if err := db.Save(&userAct).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Başarılı kayıt",
	})

}

type UptSuggest struct {
	ActionId int    `json:"action_id"`
	Action   string `json:"action"`
	Status   string `json:"status"`
}

// @Summary      Kullanıcı aksiyon güncelleme
// @Description  Aksiyon Güncelleme Kısmı
// @Tags         Action
// @Accept       json
// @Produce      json
// @Param        Suggest body UptSuggest true "status değeri için sadece 'planned,in_progress,completed,cancelled' değerlerden biri gönderilmelidir. "
// @Success      200 {object} Response "success"
// @Failure      400 {object} Response "Invalid request"
// @Router       /update-suggest [post]
func UpdateSuggest(c *gin.Context) {
	db := model.GetDB()
	user := model.User{}
	var suggest UptSuggest
	if err := c.ShouldBindJSON(&suggest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı bulunamadı.",
		})
		return
	}

	userSuggest := model.UserAction{}
	if err := db.Where("user_id=? AND id=?", user.Id, suggest.ActionId).First(&userSuggest).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if userSuggest.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Aksiyon kaydı bulunamadı.",
		})
		return
	}

	if suggest.Status != "" {
		isStatus := false
		status := []string{"planned", "in_progress", "completed", "cancelled"}
		for _, st := range status {
			if suggest.Status == st {
				isStatus = true
				break
			}
		}
		if !isStatus {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Geçersiz aksiyon status ifadesi girilmiştir.",
			})
			return
		}
		if suggest.Status != userSuggest.Status {
			userSuggest.Status = suggest.Status
		}
	}

	if suggest.Action != "" {
		if suggest.Action != userSuggest.Action {
			userSuggest.Action = suggest.Action
		}
	}

	if err := db.Save(&userSuggest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Başarılı",
	})
}

type DelSuggest struct {
	ActionId int `json:"action_id"`
}

// @Summary      Kullanıcı aksiyon silme
// @Description  Aksiyon Silme Kısmı
// @Tags         Action
// @Accept       json
// @Produce      json
// @Param        Suggest body DelSuggest true "action id gönderilmesi yeterlidir."
// @Success      200 {object} Response "success"
// @Failure      400 {object} Response "Invalid request"
// @Router       /delete-suggest [post]
func DeleteSuggest(c *gin.Context) {
	db := model.GetDB()
	user := model.User{}
	var suggest DelSuggest
	if err := c.ShouldBindJSON(&suggest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı bulunamadı.",
		})
		return
	}

	userSuggest := &model.UserAction{}
	if err := db.Where("user_id=? AND id=?", user.Id, suggest.ActionId).First(&userSuggest).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if userSuggest.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Aksiyon kaydı bulunamadı.",
		})
		return
	}
	if err := db.Delete(&model.UserAction{}, userSuggest.Id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Başarılı",
	})
}

type UserRank struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
	Rank  int     `json:"rank"`
}
type RankedUser struct {
	UserID   int
	MinScore float64
}

type RankResp struct {
	TotalScoreNumber int        `json:"total_score_number"`
	ScoreNumber      int        `json:"score_number"`
	ScoreRanking     []UserRank `json:"score_ranking"`
}

// @Summary      Kullanıcı skor sıralaması
// @Description  Score Ranking
// @Tags         Score
// @Accept       json
// @Produce      json
// @Success      200 {object} RankResp "total_score_number toplam skor tablosundaki sayıyı gösterir. score_number kullanıcının sıralamasını score_ranking ise ilk 4 sıralamayı gösterir. rank sıralamadır."
// @Failure      400 {object} Response "Invalid request"
// @Router       /score-rank [get]
func TotalRank(c *gin.Context) {
	db := model.GetDB()
	var rankedUser []RankedUser
	var userIDs []int
	user := &model.User{}
	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
	}
	subQuery := db.Model(&model.User{}).
		Select("id").
		Where("user_type = ?", user.UserType)

	lastScoreSubQuery := db.Model(&model.UserScore{}).
		Select("user_id, score, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY id DESC) as rn").
		Where("score != ?", 0).
		Where("user_id IN (?)", subQuery)

	db.Table("(?) as last_scores", lastScoreSubQuery).
		Where("rn = ?", 1).
		Select("user_id, score as min_score").
		Order("score ASC").
		Limit(3).
		Scan(&rankedUser)

	for _, rank := range rankedUser {
		userIDs = append(userIDs, rank.UserID)
	}

	var users []model.User
	if err := db.Where("id IN (?)", userIDs).Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	var response []UserRank
	for _, user := range users {
		for i, rank := range rankedUser {
			if user.Id == rank.UserID {
				response = append(response, UserRank{
					Name:  user.Username,
					Score: rank.MinScore,
					Rank:  i + 1,
				})
			}
		}
	}
	var userScore float64
	var rank int64
	var totalScore []model.UserScore
	lastScores := db.Table("(SELECT user_id, score, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY id DESC) as rn FROM user_score WHERE score != 0) as subquery").
		Where("rn = ? AND user_id IN (?)", 1, subQuery).
		Find(&totalScore)

	db.Model(&model.UserScore{}).
		Where("user_id = ?", user.Id).
		Order("id DESC").
		Limit(1).
		Pluck("score", &userScore)
	if userScore == 0 {
		rank = 0
	} else {
		db.Table("(?) as last_scores", lastScores).
			Where("user_id != ?", user.Id).
			Where("score < ?", userScore).
			Count(&rank)
		rank += 1
	}
	c.JSON(http.StatusOK, RankResp{
		TotalScoreNumber: len(totalScore),
		ScoreNumber:      int(rank),
		ScoreRanking:     response,
	})

}
