package main

import (
	"carbonfootprint/config"
	"carbonfootprint/model"
	router "carbonfootprint/routers"
	"log"

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
	model.SetDB()
	r := gin.Default()
	router.Load(r)
	r.Run(":8000")

}
