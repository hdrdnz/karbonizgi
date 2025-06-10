package controller

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

var AdminClaims jwt.MapClaims

func AdminRequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := model.GetDB()
		custom := c.GetHeader("X-Admin-Token")
		config := config.GetConfig()
		if custom == "" || custom != config.Custom.Admin {
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
		token, err := AdminValidateToken(tokenString)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				adminToken := model.AdminToken{}
				if err := db.Where("token=?", tokenString).First(&adminToken).Error; err != nil && err != gorm.ErrRecordNotFound {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"status":  "error",
						"message": "Bir hata oluştu",
					})
					return

				}
				if adminToken.Id != 0 {
					if err := db.Delete(&model.AdminToken{}, adminToken.Id).Error; err != nil {
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
		AdminClaims = claims
		c.Next()
	}
}
func AdminValidateToken(tokenString string) (*jwt.Token, error) {
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
