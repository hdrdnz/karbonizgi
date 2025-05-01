package router

import (
	"carbonfootprint/controllers"
	admin "carbonfootprint/controllers/admin"

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
	router.GET("/data", controllers.GetInfo)
	router.GET("/cal-info", controllers.GetCalInfo)

	router.GET("/comments", controllers.GetComments)
	router.GET("/suggested", controllers.GetSuggested)

	router.POST("/comp", controllers.PostComp)
	router.POST("/comp-sub", controllers.PostSubComp)

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/test", admin.GetQuestions)
		adminGroup.GET("/data", admin.GetArticle)
		adminGroup.GET("/comments", admin.GetComment)
		adminGroup.POST("/add-comment", admin.AddComment)
	}

	router.Use(controllers.RequireAuth())

	router.POST("/general-chat", controllers.GeneralChat)
	router.POST("/score", controllers.Score)
	router.GET("/user", controllers.GetUserInfo)
	router.POST("/chat", controllers.PostChat)
	router.GET("/detail-score", controllers.GetSubDetailScore)
	router.GET("/score-info", controllers.GetAllScores)
	router.POST("/add-suggest", controllers.AddSuggest)
	router.GET("/user-suggest", controllers.UserSuggest)
	router.POST("/update-suggest", controllers.UpdateSuggest)

}
