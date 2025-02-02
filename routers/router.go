package router

import (
	"carbonfootprint/controllers"
	user "carbonfootprint/controllers/user"

	_ "carbonfootprint/docs"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Load(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
	router.POST("/register", user.Register)
	router.Use(user.RequireAuth())
	router.POST("/login", user.Login)

	router.GET("/person_questions", controllers.GetPersonQues)
	router.GET("/test", controllers.Test)
	router.GET("/test2", controllers.Test2)

}
