package app

import (
	"github.com/gin-gonic/gin"
	"pmsGo/lib/config"
	"pmsGo/lib/database"
	"pmsGo/middleware"
	"pmsGo/router"
	"strconv"
)

func Run() {
	mode := gin.ReleaseMode
	if config.Config.Site.Debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	database.Init()
	server := gin.Default()
	middleware.Middleware(server)
	router.Router(server)
	server.Run(":" + strconv.Itoa(config.Config.Site.Port))
}
