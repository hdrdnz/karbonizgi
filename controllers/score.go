package controllers

import (
	"carbonfootprint/model"
	"errors"
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
	SubKey string `json:"sub_key"`
	Score  int    `json:"score"`
}

// @Summary      User Puan Kısmı
// @Description  Karbon Ayak İzi Hesaplanması
// @Tags         Score
// @Accept       json
// @Produce      json
// @Param        ScoreInfo body ScoreInfo true "/company-questions ve /person-questions endpointlerinde belli bir key value değerine göre soruları döndüm. Bu değerleri kullanarak question_name ve question_key değerlerini ekleyebilirsin."
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
				"message": err,
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
	subScores := 0
	for _, sub := range subTypes {
		questionSub := model.SubScore{}
		if sub.SubKey == "" {
			return errors.New("Soru alt başlık boş olamaz")
		}
		var subIds []int
		db.Table("question_subhead").Where("question_key=?", sub.SubKey).Pluck("id", &subIds)
		if len(subIds) == 0 {
			return errors.New("Geçersiz alt soru tipi girilmiştir.")
		}
		questionSub.QuestionSubheadId = subIds[0]
		questionSub.UserScoreId = userScore.Id
		questionSub.Score = sub.Score

		if err := db.Save(&questionSub).Error; err != nil {
			return errors.New("Bir hata oluştu.")
		}
		subScores += sub.Score
	}

	userScore.TotalScore = float64(subScores)
	if err := db.Save(&userScore).Error; err != nil {
		return errors.New("Bir hata oluştu.")

	}
	return nil
}
