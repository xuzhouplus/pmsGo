package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/lib/controller"
)

type carousel struct {
	controller.App
}

var Carousel = &carousel{}

func (controller carousel) Index(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	ctx.ShouldBind(&requestData)
	result, err := model.CarouselModel.List(0, 0, requestData["fields"], requestData["like"], requestData["order"])
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), result, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), result, "获取轮播图列表成功"))
	}
}

func (controller carousel) List(ctx *gin.Context) {
	result, err := model.CarouselModel.List(0, 0, nil, nil, nil)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), result, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), result, "获取轮播图列表成功"))
	}
}
