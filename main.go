package main

import (
	"carbonfootprint/config"
	"carbonfootprint/model"
	router "carbonfootprint/routers"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title KARBONİZGİ
// @host https://karbonizgi.leaflove.com.tr
func main() {
	_, err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatal("Config dosyası yüklenemedi:", err)
	}
	config.LoadRedis()
	config.LoadClient()
	config.GetEnv()
	model.SetDB()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","), // React dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Static("/upload", "./upload")
	router.Load(r)
	r.Run(":8000")

}
