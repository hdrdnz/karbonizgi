package controllers

import (
	"log"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası yüklenemedi!:", err)
	}
}
