package app

import (
	"github.com/gin-gonic/gin"
	"pmsGo/lib/config"
	"pmsGo/middleware"
	"pmsGo/router"
	"strconv"
)

func Start() {
	mode := gin.ReleaseMode
	if config.Config.Site.Debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	server := gin.Default()
	middleware.Middleware(server)
	router.Router(server)
	server.Run(":" + strconv.Itoa(config.Config.Site.Port))
}
