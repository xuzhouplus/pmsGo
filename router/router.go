package router

import (
	"pmsGo/controller"
	"pmsGo/lib/core"
)

func Router(engine *core.Engine) {
	//支持上传文件的访问
	engine.Static("/public", "./public")
	engine.Router(controller.Admin)
	engine.Router(controller.Nas)
	engine.Router(controller.Setting)
	engine.Router(controller.Post)
	engine.Router(controller.File)
	engine.Router(controller.Carousel)
	engine.Router(controller.CamelCase)
}
