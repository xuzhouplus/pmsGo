package middleware

import (
	"github.com/gin-gonic/gin"
	"pmsGo/lib/config"
	"pmsGo/middleware/cors"
)

func Middleware(engine *gin.Engine) {
	//根据配置进行设置跨域
	if config.Config.Site.Debug {
		engine.Use(cors.Register())
	}
}
