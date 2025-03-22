package controllers

import (
	"carbonfootprint/model"
	"net/http"
	"net/mail"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type UserRegister struct {
	Email       string `json:"email"`
	Firstname   string `json:"first_name"`
	Lastname    string `json:"last_name"`
	UserName    string `json:"user_name"`
	UserType    string `json:"user_type"`
	CompanyName string `json:"company_name"`
	Password    string `json:"password"`
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
// @Param        UserRegister body UserRegister true "user_type kısmına kullanıcı kişi ise 'person' şirket ise 'company' girmelisin. Kullanıcı bireysel ise company_name girmene gerek yok."
// @Success      200 {object} Response "success"
// @Failure      400 {object} Response "Invalid request"
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
	db.Table("user").Where("email = ?", register.Email).Count(&count)

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
			"message": "Lütfen geçerli  bir kullanıcı adı giriniz.",
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
		if register.CompanyName == "" || len(register.CompanyName) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Geçerli bir şirket ismi giriniz.",
			})
			return
		}
		user.CompanyName = register.CompanyName
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

	password, err := HashPassword(register.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}

	user.Password = password
	if err := db.Save(&user).Error; err != nil {
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
	if !CheckPasswordHash(login.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Hatalı kullanıcı adı ya da şifre1",
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

// @Summary      User Logout
// @Description  Çıkış kısmı
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} Response "User Logout successfully"
// @Router       /logout [post]
func Logout(c *gin.Context) {
	db := model.GetDB()
	userToken := &model.UserToken{}

	if err := db.Where("id=?", Claims["userId"]).First(&userToken).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	if userToken.Id != 0 {
		if err := db.Delete(&model.UserToken{}, userToken.Id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Bir hata oluştu",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Başarılı çıkış",
	})
}
