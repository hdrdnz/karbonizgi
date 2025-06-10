package controller

import (
	"carbonfootprint/controllers"
	"carbonfootprint/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type Admin struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var secretKey = []byte("./secret.key")

func AdminLogin(c *gin.Context) {
	db := model.GetDB()
	var login Admin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	admin := model.Admin{}
	if err := db.Where("email=?", login.Email).First(&admin).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if admin.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Hatalı kullanıcı adı ya da şifre",
		})
		return
	}

	if login.Password == "" || len(login.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen geçerli bir şifre giriniz.",
		})
		return
	}
	if !controllers.CheckPasswordHash(login.Password, admin.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Hatalı kullanıcı adı ya da şifre",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": admin.Id,
			"type":   "admin",
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}

	//admin token kayıt kontrolü
	adminToken := &model.AdminToken{}
	if err := db.Where("admin_id=?", admin.Id).Last(&adminToken).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if adminToken.Id == 0 {
		adminToken.AdminId = admin.Id
	}
	adminToken.Token = tokenString

	if err := db.Save(&adminToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": tokenString,
		},
	})
}
func GetAdmin(c *gin.Context) {
	db := model.GetDB()
	admin := model.Admin{}
	if err := db.Where("id=?", AdminClaims["userId"]).First(&admin).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   admin,
	})

}

func AddAdmin(c *gin.Context) {
	db := model.GetDB()
	var admin Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if admin.Name == "" && admin.LastName == "" && admin.Email == "" && admin.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bütün alanları doldurunuz.",
		})
		return
	}
	var count int64
	db.Table("admin").Where("email=?", admin.Email).Count(&count)
	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Email kayıtlıdır.",
		})
		return
	}
	password, err := controllers.HashPassword(admin.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if err := db.Save(&model.Admin{
		Name:     admin.Name,
		LastName: admin.LastName,
		Email:    admin.Email,
		Password: password,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}

func UpdateAdmin(c *gin.Context) {
	db := model.GetDB()
	var admin Admin
	adminId := c.Param("admin_id")
	recAdmin := model.Admin{}
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if err := db.Where("id=?", adminId).First(&recAdmin).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if recAdmin.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Admin bulunamadı.",
		})
		return
	}
	if admin.Name == "" && admin.LastName == "" && admin.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bütün alanları doldurunuz.",
		})
		return
	}
	recAdmin.Email = admin.Email
	recAdmin.Name = admin.Name
	recAdmin.LastName = admin.LastName
	if err := db.Save(&recAdmin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu2.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "güncelleme başarılı",
	})

}

func AdminResetPassword(c *gin.Context) {
	db := model.GetDB()
	adminId := c.Param("admin_id")
	password := c.PostForm("password")
	recAdmin := model.Admin{}
	if err := db.Where("id=?", adminId).First(&recAdmin).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	pass, err := controllers.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	recAdmin.Password = pass
	if err := db.Save(&recAdmin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "şifre başarıyla güncellendi.",
	})

}
