package main

import (
	"carbonfootprint/controllers"
	"carbonfootprint/model"
	router "carbonfootprint/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	controllers.LoadEnv()
	model.SetDB()
	r := gin.Default()
	router.Load(r)
	r.Run(":8000")

}
