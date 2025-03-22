package router

import (
	"carbonfootprint/controllers"

	_ "carbonfootprint/docs"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Load(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
	router.POST("/register", controllers.Register)
	router.GET("/person-questions", controllers.GetPersonQues)
	router.GET("/company-questions", controllers.GetCompanyQues)
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)

	router.Use(controllers.RequireAuth())
	router.POST("/score", controllers.Score)
	router.POST("/chat", controllers.PostChat)

}
