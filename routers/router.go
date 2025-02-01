package router

import (
	"carbonfootprint/controllers"

	"github.com/gin-gonic/gin"
)

func Load(router *gin.Engine) {
	router.GET("/person_questions", controllers.GetPersonQues)
	router.GET("/test", controllers.Test)

}
