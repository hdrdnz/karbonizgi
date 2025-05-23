package controllers

import (
	"carbonfootprint/model"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScoreInfo struct {
	QuestionName string        `json:"question_name"`
	QuestionSub  []QuestionSub `json:"question_key"`
}

type QuestionSub struct {
	SubKey string  `json:"sub_key"`
	Score  float64 `json:"score"`
}

type ScoreResp struct {
	ScoreId int    `json:"score_id"`
	Score   string `json:"score"`
}

type ResponseScore struct {
	Data    ScoreResp `json:"data"`
	Message string    `json:"message"`
	Status  string    `json:"status"`
}

// @Summary      User Puan Kısmı
// @Description  Karbon Ayak İzi Hesaplanması
// @Tags         Score
// @Accept       json
// @Produce      json
// @Param        ScoreInfo body ScoreInfo true "/company-questions ve /person-questions endpointlerinde belli bir key value değerleri bulunamktadır. Bu değerleri kullanarak question_name ve question_key değerleri eklenir."
// @Success      200 {object} ResponseScore "data kısmında kullanıcının karbon ayak izi değeri döner.Chat kısmında bu kısmı kullanacaksın."
// @Failure      400 {object} Response "Invalid request"
// @Router       /score [post]
func Score(c *gin.Context) {
	fmt.Println("girdiii")
	db := model.GetDB()
	var scoreInf []ScoreInfo
	if err := c.ShouldBindJSON(&scoreInf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	userId := Claims["userId"].(float64)
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı bulunamadı.",
		})
		return
	}
	user := model.User{}
	if err := db.Where("id=?", userId).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	score := model.UserScore{}
	if err := db.Where("user_id=?", user.Id).Last(&score).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	// nowDate := time.Now()
	// newDate := score.CreatedAt.AddDate(0, 0, 7)
	// if nowDate.Before(newDate) {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  "error",
	// 		"message": "Skor güncellemesi bir hafta ara ile yapılmaktadır.",
	// 	})
	// 	return
	// }

	//toplam skorun eklenmesi
	totalScore := 0.0
	for _, info := range scoreInf {
		for _, sub := range info.QuestionSub {
			totalScore += sub.Score
		}
	}
	totalScore = math.Round((totalScore/1000.0)*10) / 10
	userScore := model.UserScore{}
	userScore.Score = totalScore
	userScore.UserId = int(userId)
	if err := db.Save(&userScore).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	fmt.Println("başarıyla eklendi")
	for _, info := range scoreInf {
		detailScore := &model.UserDetailScore{}
		detailScore.UserId = int(userId)
		detailScore.UserScoreId = userScore.Id
		question := model.QuestionTypes{}
		if info.QuestionName == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Soru tipi kısmı boş olamaz.",
			})
			return
		}
		if err := db.Where("question_key=?", info.QuestionName).First(&question).Error; err != nil && err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Soru tipi kısmı boş olamaz.",
			})
			return
		}

		if question.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Geçersiz soru tipi girilmiştir.",
			})
			return
		}

		if question.QuestionType != user.UserType {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Geçersiz soru tipi girilmiştir.",
			})
			return
		}

		detailScore.QuestionTypesId = question.Id
		if err := db.Save(&detailScore).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Bir hata oluştu.",
			})
			return
		}

		//alt başlıkların eklenmesi
		err := subSccore(info.QuestionSub, detailScore)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		// totalScore += detailScore.TotalScore
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sonuçlar başarılı bir şekilde eklenmiştir.",
		"data": ScoreResp{
			ScoreId: userScore.Id,
			Score:   strconv.FormatFloat(totalScore, 'f', 1, 64),
		},
	})
}

func subSccore(subTypes []QuestionSub, userScore *model.UserDetailScore) error {
	db := model.GetDB()
	subScores := 0.0
	for _, sub := range subTypes {
		questionSub := model.SubScore{}
		if sub.SubKey == "" {
			return errors.New("Soru alt başlık boş olamaz")
		}
		var subIds []int
		db.Table("question_subhead").Where("question_key=?", sub.SubKey).Pluck("id", &subIds)
		if len(subIds) == 0 {
			return errors.New("Geçersiz alt soru tipi girilmiştir")
		}
		questionSub.QuestionSubheadId = subIds[0]
		questionSub.UserDetailScoreId = userScore.Id
		questionSub.Score = sub.Score
		if err := db.Save(&questionSub).Error; err != nil {
			return errors.New("Bir hata oluştu")
		}
		subScores += sub.Score
	}

	userScore.TotalScore = float64(subScores)
	if err := db.Save(&userScore).Error; err != nil {
		return errors.New("Bir hata oluştu")

	}
	return nil
}

type DetailScore struct {
	Id              int     `json:"questionSub_Id"`
	QuestionType    string  `json:"questionType"`
	QuestionSubType string  `json:"questionSubType"`
	SubScore        float64 `json:"subScore"`
}

type DetailResp struct {
	Data []DetailScore `json:"data"`
}

// @Description  Detay skor bilgileri
// @Tags         Score
// @Accept       json
// @Produce      json
// @Param         score_id query string true "skor id değeri girilir."
// @Success      200 {object} DetailResp "Kullanıcıya ait soru alt başlıklarına göre questiontype:temel soru başlığını ,QuestionSubType : alt soru başlığını ve SubScore ise alt başlığa ait değeri içerir.Bu kısımlar chat kısmı için kullanılır."
// @Router       /detail-score [get]
func GetSubDetailScore(c *gin.Context) {
	user := model.User{}
	db := model.GetDB()
	scoreId := c.Query("score_id")
	score := model.UserScore{}
	if scoreId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Skor id boş bırakılamaz.",
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
	if err := db.Where("id=? and user_id=?", scoreId, user.Id).First(&score).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if score.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Skor bulunamadı.",
		})
		return
	}
	var userIds []int
	if err := db.Table("user_detail_score").Where("user_id=? and user_score_id=?", user.Id, score.Id).Pluck("id", &userIds).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	fmt.Println("userIDS:", userIds)

	subScores := []model.SubScore{}
	if err := db.Where("user_detail_score_id IN (?)", userIds).Preload("QuestionSubhead").Preload("UserDetailScore").Preload("QuestionSubhead.QuestionTypes").Find(&subScores).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	var allScores []DetailScore
	for _, item := range subScores {
		score := DetailScore{
			Id:              item.Id,
			QuestionType:    item.QuestionSubhead.QuestionTypes.QuestionKey,
			QuestionSubType: item.QuestionSubhead.QuestionKey,
			SubScore:        item.Score,
		}
		allScores = append(allScores, score)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   allScores,
	})
}

type Scores struct {
	ScoreId   int     `json:"score_id"`
	Score     float64 `json:"score"`
	ScoreDate string  `json:"score_date"`
}

type DataScore struct {
	Data   []AllScores `json:"data"`
	Status string      `json:"status"`
}
type AllScores struct {
	Scores    []Scores `json:"scores"`
	LastScore Scores   `json:"last-score"`
}

// @Summary      Kullanıcı skor bilgileri
// @Description  Kullanıcı skor bilgileri
// @Tags         Score
// @Accept       json
// @Produce      json
// @Success      200 {object} DataScore "Kullanıcıya ait skor puanını ve skor tarihini geri döndürür."
// @Router       /score-info [get]
func GetAllScores(c *gin.Context) {
	db := model.GetDB()
	scores := []model.UserScore{}
	if err := db.Where("user_id=?", Claims["userId"]).Find(&scores).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}

	if len(scores) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcıya ait skor bulunmamaktadır.",
		})
		return
	}
	allScores := []Scores{}
	id := 0
	var lastScore Scores

	for _, score := range scores {
		if score.Id > id {
			lastScore.ScoreId = score.Id
			lastScore.Score = score.Score
			lastScore.ScoreDate = score.CreatedAt.Format("02-01-2006")
		}
		allScores = append(allScores, Scores{
			ScoreId:   score.Id,
			Score:     score.Score,
			ScoreDate: score.CreatedAt.Format("02-01-2006"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": AllScores{
			Scores:    allScores,
			LastScore: lastScore,
		},
	})
}

// func RankScore(c *gin.Context) {
// 	db := model.GetDB()
// 	var lowestScores []model.UserScore
// 	db.Order("score ASC").Limit(4).Find(&lowestScores)
// 	userId := Claims["userId"]

// }

func CompTest(c *gin.Context) {
	var scoreInf []ScoreInfo
	if err := c.ShouldBindJSON(&scoreInf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	totalScore := 0.0
	for _, info := range scoreInf {
		for _, sub := range info.QuestionSub {
			totalScore += sub.Score
		}
	}
	totalScore = math.Round((totalScore/1000.0)*10) / 10
	c.JSON(http.StatusOK, gin.H{
		"total": totalScore,
	})
}
