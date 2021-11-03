package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
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

// Init 初始化服务器
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

// Static 注册静态资源访问路由
func (e *Engine) Static(prefix, path string) {
	e.server.Static(prefix, path)
}

// Use 注册中间件
func (e *Engine) Use(handlerFunc ...gin.HandlerFunc) {
	e.server.Use(handlerFunc...)
}

// Router 注册控制器
func (e *Engine) Router(appInterface controller.AppInterface) {
	//解析控制器，获取控制器名称、方法、方法请求方式和方法登录认证配置
	resolve := controller.NewResolve(appInterface)
	//注册路由分组
	group := e.server.Group("/" + resolve.GetControllerName())
	//添加登录认证配置
	auth.Add(resolve)
	//遍历方法，注册子路由
	for _, action := range resolve.GetActions() {
		//遍历方法请求方式
		for _, verb := range action.Verbs {
			var actionName string
			if action.Name == "index" {
				actionName = ""
			} else {
				actionName = "/" + action.Name
			}
			//如果是*注册所有请求方法，经过特殊处理，*一定出现在第一个位置
			if verb == controller.Any {
				group.Any(actionName, resolve.Handle(action.Name))
				break
			}
			//注册对应请求方式的路由
			group.Handle(strings.ToUpper(verb), actionName, resolve.Handle(action.Name))
		}
	}
}

// Start 启动服务
func (e *Engine) Start(listen string) error {
	//监听端口
	err := e.server.Run(listen)
	if err != nil {
		return err
	}
	return nil
}

func (e Engine) Pprof(listen string) {
	go func() {
		http.ListenAndServe(listen, nil)
	}()
}
