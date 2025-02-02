package controllers

import (
	"carbonfootprint/controllers"
	"carbonfootprint/model"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var secretKey = []byte("./secret.key")
var Claims jwt.Claims

type UserRegister struct {
	Email     string `json:"email"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	UserName  string `json:"user_name"`
	UserType  string `json:"user_type"`
	Password  string `json:"password"`
}

type UserLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// @Summary      User Registration
// @Description  Kayıt kısmı
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        UserRegister body UserRegister true "User registration details"
// @Success      200 {object} Response "User registered successfully"
// @Failure      400 {object} Response "Invalid request"
// @Failure      500 {object} Response "Internal server error"
// @Router       /register [post]
func Register(c *gin.Context) {
	db := model.GetDB()
	var register UserRegister
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}

	user := &model.User{}
	_, err := mail.ParseAddress(register.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen geçerli bir e-posta adresi giriniz.",
		})
		return
	}
	var count int64
	db.Table("user").Where("email=?", register.Email).Count(&count)

	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bu mail ile kayıtlı kullanıcı bulunmaktadır.",
		})
		return
	}

	user.Email = register.Email

	if register.Firstname == "" || len(register.Firstname) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen geçerli bir isim giriniz.",
		})
		return
	}

	user.Firstname = register.Firstname

	if register.Lastname == "" || len(register.Firstname) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen geçerli  bir  soyisim giriniz.",
		})
		return
	}

	user.Lastname = register.Lastname

	if register.UserName == "" || len(register.UserName) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen geçerli  bir  soyisim giriniz.",
		})
		return
	}

	db.Table("user").Where("username=?", register.UserName).Count(&count)

	if count != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bu kullanıcı adı kullanılıyor.Başka bir kullanıcı adı giriniz.",
		})
		return
	}
	user.Username = register.UserName

	if register.UserType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı tipini seçiniz.",
		})
		return
	}

	if register.UserType == "person" {
		user.UserType = "person"
	} else if register.UserType == "company" {
		user.UserType = "company"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Kullanıcı tipi hatalı",
		})
		return
	}

	if len(register.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Şifre en az 6 karakter olmalıdır.",
		})
		return
	}

	upperCase := regexp.MustCompile(`[A-Z]`)
	if !upperCase.MatchString(register.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "şifre en az bir büyük harf içermelidir",
		})
		return
	}

	password, err := controllers.HashPassword(register.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}

	user.Password = password
	if err := db.Save(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kullanıcı başarıyla oluşturuldu.",
	})
}

// @Summary      User Login
// @Description  Giriş kısmı
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        UserLogin body UserLogin true "User login details"
// @Success      200 {object} Response "User login successfully"
// @Failure      400 {object} Response "Invalid request"
// @Failure      500 {object} Response "Internal server error"
// @Router       /login [post]
func Login(c *gin.Context) {
	db := model.GetDB()
	var login UserLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	user := &model.User{}
	if login.UserName == "" || len(login.UserName) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen geçerli bir kullanıcı adı giriniz.",
		})
		return
	}
	if err := db.Where("username=?", login.UserName).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}

	if user.Id == 0 {
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

	if !controllers.CheckPasswordHash(login.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Hatalı kullanıcı adı ya da şifre",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":   user.Id,
			"userName": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}

	//kullanıcı token kayıt kontrolü
	userToken := &model.UserToken{}
	if err := db.Where("id=?", user.Id).Last(&userToken).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	if userToken.Id == 0 {
		userToken.UserId = user.Id
	}
	userToken.Token = tokenString

	if err := db.Save(&userToken).Error; err != nil {
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

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		custom := c.GetHeader("Custom-Header")
		if custom == "" || custom != os.Getenv("Custom-Header") {
			c.JSON(http.StatusBadRequest, nil)
			c.Abort()
		}
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Yetkisi erişim",
			})
			c.Abort()
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		token, err := ValidateToken(tokenString)
		if err != nil {
			c.Abort()
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token hatası",
			})
			c.Abort()
		}
		Claims = claims
		c.Next()
	}
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("geçersiz imzalama yöntemi")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
