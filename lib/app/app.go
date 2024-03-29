package app

import (
	"fmt"
	"pmsGo/lib/config"
	"pmsGo/lib/core"
	"pmsGo/lib/middleware/auth"
	"pmsGo/lib/middleware/cors"
	"pmsGo/lib/middleware/session"
	"pmsGo/lib/sync"
	"pmsGo/router"
)

// Bootstrap yaml配置、数据库连接、数据库配置在各自包的init方法中初始化
func Bootstrap() {
	core.App.Init(config.Config.Site.Debug)
	if config.Config.Site.Debug {
		core.App.Pprof(":8081")
	}
	//根据配置进行设置跨域
	if config.Config.Site.AllowCrossDomain {
		core.App.Use(cors.Register())
	}
	//注册session中间件
	core.App.Use(session.Register())
	//注册登录验证中间件
	core.App.Use(auth.Auth())
	//注册路由
	router.Router(core.App)
	//启动异步任务队列
	sync.Run()
	//启动web服务
	err := core.App.Start(config.Config.Site.Listen)
	if err != nil {
		fmt.Printf("%err\n", err)
		return
	}
}
