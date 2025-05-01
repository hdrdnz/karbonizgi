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

// @Summary      User Puan Kısmı
// @Description  Karbon Ayak İzi Hesaplanması
// @Tags         Score
// @Accept       json
// @Produce      json
// @Param        ScoreInfo body ScoreInfo true "/company-questions ve /person-questions endpointlerinde belli bir key value değerleri bulunamktadır. Bu değerleri kullanarak question_name ve question_key değerleri eklenir."
// @Success      200 {object} Response "data kısmında kullanıcının karbon ayak izi değeri döner.Chat kısmında bu kısmı kullanacaksın."
// @Failure      400 {object} Response "Invalid request"
// @Router       /score [post]
func Score(c *gin.Context) {
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
	totalScore := 0.0
	for _, info := range scoreInf {
		score := &model.UserDetailScore{}
		score.UserId = int(userId)
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

		score.QuestionTypesId = question.Id
		if err := db.Save(&score).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Bir hata oluştu.",
			})
			return
		}

		//alt başlıkların eklenmesi
		err := subSccore(info.QuestionSub, score)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		totalScore += score.TotalScore
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
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sonuçlar başarılı bir şekilde eklenmiştir.",
		"data":    strconv.FormatFloat(totalScore, 'f', 1, 64),
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
// @Success      200 {object} DetailResp "Kullanıcıya ait soru alt başlıklarına göre questiontype:temel soru başlığını ,QuestionSubType : alt soru başlığını ve SubScore ise alt başlığa ait değeri içerir.Bu kısımlar chat kısmı için kullanılır."
// @Router       /detail-score [get]
func GetSubDetailScore(c *gin.Context) {
	user := model.User{}
	db := model.GetDB()
	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	var userIds []int
	if err := db.Table("user_detail_score").Where("user_id=?", user.Id).Pluck("id", &userIds).Error; err != nil && err != gorm.ErrRecordNotFound {
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

type AllScores struct {
	Score     float64 `json:"score"`
	ScoreDate string  `json:"score_date"`
}

type DataScore struct {
	Data   []AllScores `json:"data"`
	Status string      `json:"status"`
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
	allScores := []AllScores{}

	for _, score := range scores {
		allScores = append(allScores, AllScores{
			Score:     score.Score,
			ScoreDate: score.CreatedAt.Format("02-01-2006"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   allScores,
	})
}
