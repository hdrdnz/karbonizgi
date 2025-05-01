package controllers

import (
	"carbonfootprint/config"
	"carbonfootprint/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

var Claims jwt.MapClaims
var secretKey = []byte("./secret.key")

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("secretKey:", secretKey)
		db := model.GetDB()
		custom := c.GetHeader("X-Custom-Token")
		config := config.GetConfig()
		fmt.Println("onfig.Custom.Header:", config.Custom.Header)

		if custom == "" || custom != config.Custom.Header {
			c.AbortWithStatusJSON(http.StatusBadRequest, nil)
			return
		}
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Yetkisiz erişim",
			})
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		token, err := ValidateToken(tokenString)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				userToken := model.UserToken{}
				if err := db.Where("token=?", tokenString).First(&userToken).Error; err != nil && err != gorm.ErrRecordNotFound {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"status":  "error",
						"message": "Bir hata oluştu",
					})
					return

				}
				if userToken.Id != 0 {
					if err := db.Delete(&model.UserToken{}, userToken.Id).Error; err != nil {
						c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
							"status":  "error",
							"message": "Bir hata oluştu",
						})
						return

					}
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  "error",
					"message": "Token süresi dolmuştur.Lütfen tekrardan giriş yapınız.",
				})
				return
			}
		}
		if token == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token hatalı",
			})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token hatalı",
			})
			return
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
