package router

import (
	"github.com/gin-gonic/gin"
	"pmsGo/app/controller"
)

func Router(engine *gin.Engine) {
	admin := engine.Group("/admin")
	admin.Handle("POST", "/login", controller.Admin.Login)

	setting := engine.Group("/setting")
	setting.GET("", controller.Setting.Index)

	carousel := engine.Group("/carousel")
	carousel.GET("", controller.Carousel.Index)

	post := engine.Group("/post")
	post.GET("", controller.Post.Index)
}
