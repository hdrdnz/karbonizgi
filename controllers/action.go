package controllers

import (
	"carbonfootprint/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Suggest struct {
	Action string `json:"action"`
	Status string `json:"status"`
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

type UptSuggest struct {
	ActionId int    `json:"actionId"`
	Action   string `json:"action"`
	Status   string `json:"status"`
}

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

	fmt.Println("User:", user.Id)

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
			"message": "Kullanıcıya ait aksiyon kaydı bulunamadı.",
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
	c.JSON(http.StatusOK, "Başarılı")
}
