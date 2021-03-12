package router

import (
	"github.com/gin-gonic/gin"
	"pmsGo/app/controller"
)

func Router(engine *gin.Engine) {
	admin := engine.Group("/admin")
	admin.Handle("POST", "/login", controller.Admin.Login)
}
