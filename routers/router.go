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
	router.GET("key-translation", controllers.KeyTranslation)

	router.GET("/comments", controllers.GetComments)
	router.GET("/suggested", controllers.GetSuggested)

	router.POST("/comp", controllers.PostComp)
	router.POST("/comp-sub", controllers.PostSubComp)

	adminGroup := router.Group("/admin")
	adminGroup.POST("/login", admin.AdminLogin)
	adminGroup.Use(admin.AdminRequireAuth())
	{
		adminGroup.GET("/test", admin.GetQuestions)
		adminGroup.POST("/delete-test", admin.DeleteQuestion)
		adminGroup.POST("/add-test/:type/:category", admin.AddQuestion)
		adminGroup.GET("/data", admin.GetArticle)
		adminGroup.POST("/add-data", admin.AddArticle)
		adminGroup.POST("/delete-data", admin.DeleteArticle)
		adminGroup.GET("/comments", admin.GetComment)
		adminGroup.POST("/add-comment", admin.AddComment)
		adminGroup.POST("/delete-comment", admin.DeleteComment)
		adminGroup.GET("/users", admin.GetUser)
		adminGroup.GET("", admin.GetAdmin)
		adminGroup.POST("/add-admin", admin.AddAdmin)
		adminGroup.POST("/update-admin/:admin_id", admin.UpdateAdmin)
		adminGroup.POST("/reset-admin/:admin_id", admin.AdminResetPassword)
		adminGroup.GET("/total", admin.Total)
		adminGroup.POST("/update-user", admin.UpdateUser)
		adminGroup.POST("/add-user", admin.AddUser)
		adminGroup.POST("/delete-user", admin.DeleteUser)
		adminGroup.POST("/reset-password", admin.ResetPassword)

	}

	router.Use(controllers.RequireAuth())
	router.POST("/general-chat", controllers.GeneralChat)
	router.POST("/score", controllers.Score)
	router.POST("/score-test", controllers.CompTest)
	router.GET("/user", controllers.GetUserInfo)
	router.POST("/chat", controllers.PostChat)
	router.GET("/detail-score", controllers.GetSubDetailScore)
	router.GET("/score-info", controllers.GetAllScores)
	router.POST("/add-suggest", controllers.AddSuggest)
	router.GET("/user-suggest", controllers.UserSuggest)
	router.POST("/update-suggest", controllers.UpdateSuggest)
	router.POST("/delete-suggest", controllers.DeleteSuggest)
	router.GET("/score-rank", controllers.TotalRank)

}
