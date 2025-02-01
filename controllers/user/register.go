package controllers

import (
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
)

type UserRegister struct {
	Email     string `json:"email"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func Register(c *gin.Context) {
	var register UserRegister
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Bir hata oluştu",
		})
		return
	}

	//user := &model.User{}
	_, err := mail.ParseAddress(register.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen geçerli bir e-posta adresi giriniz.",
		})
		return
	}

	if register.Firstname == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Lütfen bir isim giriniz.",
		})
		return
	}

}
