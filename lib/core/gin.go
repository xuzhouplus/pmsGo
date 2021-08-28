package core

import (
	"github.com/gin-gonic/gin"
	"pmsGo/lib/controller"
	"pmsGo/lib/middleware/auth"
	"strings"
)

type Engine struct {
	server *gin.Engine
}

var App *Engine

func init() {
	App = &Engine{}
}

func (e *Engine) Init(debug bool) {
	//设置启动模式
	mode := gin.ReleaseMode
	if debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	//创建gin服务实例
	e.server = gin.Default()
}

func (e *Engine) Static(prefix, path string) {
	e.server.Static(prefix, path)
}

func (e *Engine) Use(handlerFunc ...gin.HandlerFunc) {
	e.server.Use(handlerFunc...)
}

func (e *Engine) Router(appInterface controller.AppInterface) {
	resolve := controller.NewResolve(appInterface)
	group := e.server.Group("/" + resolve.GetControllerName())
	auth.Add(resolve)
	for _, action := range resolve.GetActions() {
		for _, verb := range action.Verbs {
			var actionName string
			if action.Name == "index" {
				actionName = ""
			} else {
				actionName = "/" + action.Name
			}
			if verb == controller.Any {
				group.Any(actionName, resolve.Handle(action.Name))
				break
			}
			group.Handle(strings.ToUpper(verb), actionName, resolve.Handle(action.Name))
		}
	}
}

func (e *Engine) Start(listen string) error {
	//监听端口
	err := e.server.Run(listen)
	if err != nil {
		return err
	}
	return nil
}
