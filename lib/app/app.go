package app

import (
	"github.com/gin-gonic/gin"
	"pmsGo/lib/config"
	"pmsGo/router"
	"strconv"
)

func Run() {
	mode := gin.ReleaseMode
	if config.Config.Site.Debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	server := gin.Default()
	router.Router(server)
	server.Run(":" + strconv.Itoa(config.Config.Site.Port))
}
