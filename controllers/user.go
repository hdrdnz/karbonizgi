package controllers

import (
	"carbonfootprint/model"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserInfo struct {
	Id          int
	Email       string
	Firstname   string
	Lastname    string
	Username    string
	UserType    string
	CompanyName string
}

// @Description  Kullanıcı bilgileri kısmı
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200 {object} UserInfo
// @Router       /user [get]
func GetUserInfo(c *gin.Context) {
	db := model.GetDB()
	user := model.User{}
	if err := db.Where("id=?", Claims["userId"]).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	var userInfo UserInfo
	byteResp, _ := json.Marshal(user)
	if err := json.Unmarshal(byteResp, &userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	c.JSON(http.StatusOK, userInfo)
}
