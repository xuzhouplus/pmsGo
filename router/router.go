package router

import (
	"github.com/gin-gonic/gin"
	"pmsGo/app/controller"
	"pmsGo/middleware/auth"
	"pmsGo/middleware/session"
)

func Router(engine *gin.Engine) {
	engine.Static("/public", "./public")
	engine.Use(session.Register())
	admin := engine.Group("/admin")
	{
		admin.POST("/login", auth.Register(), controller.Admin.Login)
		admin.POST("/auth", auth.Register(), controller.Admin.Auth)
		admin.POST("/logout", auth.Register(), controller.Admin.Logout)
		admin.GET("/profile", auth.Register(), controller.Admin.Profile)
		admin.POST("/profile", auth.Register(), controller.Admin.Profile)
		admin.GET("/connects", auth.Register(), controller.Admin.Connects)
		admin.GET("/authorize", auth.Register(), controller.Admin.AuthorizeUrl)
		admin.GET("/authorize-user", auth.Register(), controller.Admin.AuthorizeUser)
	}

	setting := engine.Group("/setting")
	{
		setting.GET("", controller.Setting.Index)
		setting.GET("/site", auth.Register(), controller.Setting.Site)
		setting.POST("/site", auth.Register(), controller.Setting.Save)
		setting.GET("/carousel", auth.Register(), controller.Setting.Carousel)
		setting.POST("/carousel", auth.Register(), controller.Setting.Save)
		setting.GET("/alipay", auth.Register(), controller.Setting.Alipay)
		setting.POST("/alipay", auth.Register(), controller.Setting.Save)
		setting.GET("/baidu", auth.Register(), controller.Setting.Baidu)
		setting.POST("/baidu", auth.Register(), controller.Setting.Save)
		setting.GET("/facebook", auth.Register(), controller.Setting.Facebook)
		setting.POST("/facebook", auth.Register(), controller.Setting.Save)
		setting.GET("/github", auth.Register(), controller.Setting.Github)
		setting.POST("/github", auth.Register(), controller.Setting.Save)
		setting.GET("/google", auth.Register(), controller.Setting.Google)
		setting.POST("/google", auth.Register(), controller.Setting.Save)
		setting.GET("/line", auth.Register(), controller.Setting.Line)
		setting.POST("/line", auth.Register(), controller.Setting.Save)
		setting.GET("/qq", auth.Register(), controller.Setting.Qq)
		setting.POST("/qq", auth.Register(), controller.Setting.Save)
		setting.GET("/twitter", auth.Register(), controller.Setting.Twitter)
		setting.POST("/twitter", auth.Register(), controller.Setting.Save)
		setting.GET("/wechat", auth.Register(), controller.Setting.Wechat)
		setting.POST("/wechat", auth.Register(), controller.Setting.Save)
		setting.GET("/weibo", auth.Register(), controller.Setting.Weibo)
		setting.POST("/weibo", auth.Register(), controller.Setting.Save)
	}

	carousel := engine.Group("/carousel")
	{
		carousel.GET("", controller.Carousel.Index)
		carousel.GET("/list", auth.Register(), controller.Carousel.List)
		carousel.POST("/create", auth.Register(), controller.Carousel.Create)
		carousel.POST("/delete", auth.Register(), controller.Carousel.Delete)
	}

	post := engine.Group("/post")
	{
		post.GET("", controller.Post.Index)
		post.GET("/list", auth.Register(), controller.Post.List)
		post.POST("/save", auth.Register(), controller.Post.Save)
		post.POST("/delete", auth.Register(), controller.Post.Delete)
		post.POST("/toggle-status", auth.Register(), controller.Post.ToggleStatus)
		post.GET("/detail", controller.Post.Detail)
	}

	file := engine.Group("/file")
	{
		file.GET("", auth.Register(), controller.File.Index)
		file.POST("/upload", auth.Register(), controller.File.Upload)
		file.POST("/delete", auth.Register(), controller.File.Delete)
	}

}
