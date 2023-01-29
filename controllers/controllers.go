package controllers

import (
	"BillSystem/views"
	"github.com/gin-gonic/gin"
	"BillSystem/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func InitController() *gin.Engine{
	controller := gin.Default()
	controller.MaxMultipartMemory = 8 << 20
	store, _ := redis.NewStore(10, "tcp", "127.0.0.1:6379", "", []byte("secret"))
	controller.Use(sessions.Sessions("mysession", store))
	controller.Use(middlewares.Cors())
	controller.POST("/login", views.Login)
	controller.GET("/logout", views.Logout)

	v1 := controller.Group("/api/v1")
	v1.Use(middlewares.AuthRequired)
	{
		//v1.GET("/", views.TestRoot)
		//v1.GET("/test/*id", views.TestInfo)
		//v1.POST("/post/", views.TestPost)
		//v1.POST("/upload", views.TestUpload)
		v1.GET("/alyBill", views.AlyBill)
		v1.GET("/alyApportionBill", views.AlyApportionBill)
		v1.GET("/alyUnApportionBill", views.AlyUnApportionBill)
		v1.GET("/centerDepart", views.CenterDepart)
	}

	private := controller.Group("/private")
	private.Use(middlewares.AuthRequired)
	{
		private.GET("/me", views.Me)
		private.GET("/status", views.Status)
	}
	return controller
}
