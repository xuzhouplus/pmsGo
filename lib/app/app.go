package app

import (
	"github.com/gin-gonic/gin"
	"pmsGo/lib/config"
	"pmsGo/router"
	"strconv"
)

// Run yaml配置、数据库连接、数据库配置在各自包的init方法中初始化
func Run() {
	//设置启动模式
	mode := gin.ReleaseMode
	if config.Config.Site.Debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	//创建gin服务实例
	server := gin.Default()
	//注册路由
	router.Router(server)
	//监听端口
	server.Run(":" + strconv.Itoa(config.Config.Site.Port))
}
