package app

import (
	"pmsGo/lib/config"
	"pmsGo/lib/core"
	"pmsGo/lib/middleware/auth"
	"pmsGo/lib/middleware/cors"
	"pmsGo/lib/middleware/session"
	"pmsGo/router"
)

// Run yaml配置、数据库连接、数据库配置在各自包的init方法中初始化
func Bootstrap() {
	core.App.Init(config.Config.Site.Debug)
	//根据配置进行设置跨域
	if config.Config.Site.Debug {
		core.App.Use(cors.Register())
	}
	core.App.Use(session.Register())
	core.App.Use(auth.Auth())
	router.Router(core.App)
	core.App.Start(config.Config.Site.Listen)
}
