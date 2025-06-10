package controller

import (
	"carbonfootprint/controllers"
	"carbonfootprint/model"
	"encoding/json"
	"io"
	"net/http"
	"net/mail"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUser(c *gin.Context) {
	db := model.GetDB()
	userType := c.Query("user_type")
	if userType != "person" && userType != "company" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz user tipi girilmiştir.",
		})
		return
	}

	var users []model.User
	if err := db.Where("user_type=?", userType).Order("id desc").Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

type UserInfo struct {
	UserId    int
	Email     string
	UserName  string
	FirstName string
	LastName  string
	UserType  string
}

func UpdateUser(c *gin.Context) {
	db := model.GetDB()
	var updateUser UserInfo
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	user := model.User{}
	if err := db.Where("id=?", updateUser.UserId).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
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
	if updateUser.Email != "" && user.Email != updateUser.Email {
		var count int64
		db.Table("user").Where("email=?", updateUser.Email).Count(&count)
		if count != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Email kullanılmaktadır.",
			})
			return
		}
		user.Email = updateUser.Email
	}
	if updateUser.UserName != "" && user.Username != updateUser.UserName {
		var count int64
		db.Table("user").Where("username=?", updateUser.UserName).Count(&count)
		if count != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Kullanıcı adı bulunmaktadır.",
			})
			return
		}
		user.Username = updateUser.UserName

	}
	if updateUser.FirstName != "" && user.Firstname != updateUser.FirstName {
		user.Firstname = updateUser.FirstName
	}
	if updateUser.LastName != "" && user.Lastname != updateUser.LastName {
		user.Lastname = updateUser.LastName
	}
	if updateUser.UserType == "" && updateUser.UserType != "person" && updateUser.UserType != "company" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz kullanıcı tipi",
		})
		return
	}
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Geçersiz kullanıcı tipi",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Başarılı",
		"status":  "success",
	})

}
func AddUser(c *gin.Context) {
	db := model.GetDB()
	var register controllers.UserRegister
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
	if !controllers.ContainsUpper(register.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Şifre en az bir büyük harf içermelidir.",
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
func DeleteUser(c *gin.Context) {
	db := model.GetDB()
	user := model.User{}
	userId := c.PostForm("userId")
	if err := db.Where("id=?", userId).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	var scoreIds []uint
	var detailIds []uint
	db.Table("user_score").Where("user_id=?", userId).Pluck("id", &scoreIds)
	if len(scoreIds) != 0 {
		db.Table("user_detail_score").Where("user_score_id IN (?)", scoreIds).Pluck("id", &detailIds)
		if err := db.Where("id IN ?", scoreIds).
			Delete(&model.UserScore{}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Bir hata oluştu.",
			})
			return
		}
	}
	if len(detailIds) != 0 {
		if err := db.Where("user_detail_score_id IN ?", detailIds).
			Delete(&model.SubScore{}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Bir hata oluştu.",
			})
			return
		}
		if err := db.Where("id IN ?", detailIds).
			Delete(&model.UserDetailScore{}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Bir hata oluştu.",
			})
			return
		}
	}
	if err := db.Where("user_id=?", userId).
		Delete(&model.UserAction{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
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
	if err := db.Delete(&model.User{}, user.Id).Error; err != nil {
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

func ResetPassword(c *gin.Context) {
	db := model.GetDB()
	user := model.User{}
	userId := c.PostForm("userId")
	if err := db.Where("id=?", userId).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	password := c.PostForm("password")
	if len(password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "şifre uzunluğu 6 ve üzeri olmalıdır.",
		})
		return
	}
	if !controllers.ContainsUpper(password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Şifre en az bir büyük harf içermelidir.",
		})
		return
	}

	password, err := controllers.HashPassword(password)
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
		"message": "Kullanıcı şifresi güncellendi.",
	})
}

type TotalResp struct {
	User    int64
	Makale  int
	Person  int
	Company int
}

func Total(c *gin.Context) {
	db := model.GetDB()
	//kullanıcı
	var total int64
	db.Table("user").Count(&total)
	//makale
	file, _ := os.Open("./data/data.json")
	byteFile, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.",
		})
		return
	}
	var ques []map[string]interface{}
	if err := json.Unmarshal(byteFile, &ques); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu2.",
		})
		return
	}
	//bireysel test
	file, _ = os.Open("./data/person.json")
	byteFile, err = io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.1",
		})
		return
	}
	person := make(map[string]interface{})
	if err := json.Unmarshal(byteFile, &person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu2.",
		})
		return
	}
	personCount := 0
	for _, p := range person {
		personCount += len(p.([]interface{}))
	}
	//kurumsal test
	file, _ = os.Open("./data/person.json")
	byteFile, err = io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu.1",
		})
		return
	}
	company := make(map[string]interface{})
	if err := json.Unmarshal(byteFile, &company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu2.",
		})
		return
	}
	companyCount := 0
	for _, p := range person {
		companyCount += len(p.([]interface{}))
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": TotalResp{
			User:    total,
			Makale:  len(ques),
			Person:  personCount,
			Company: companyCount,
		},
	})
}
