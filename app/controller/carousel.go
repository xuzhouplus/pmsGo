package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/controller"
	"pmsGo/lib/helper/image"
	"strconv"
)

type carousel struct {
	controller.App
}

var Carousel = &carousel{}

func (controller carousel) Index(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	err := ctx.ShouldBind(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
		return
	}
	result, err := service.CarouselService.List(0, 0, requestData["fields"], requestData["like"], requestData["order"])
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), result, "获取轮播图列表成功"))
	}
}

func (controller carousel) List(ctx *gin.Context) {
	list, err := service.CarouselService.List(0, 0, nil, nil, nil)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
	} else {
		data := make(map[string]interface{})
		data["list"] = list
		carouselLimit, _ := strconv.Atoi(service.SettingService.GetSetting(model.SettingKeyCarouselLimit))
		fmt.Println(carouselLimit)
		data["limit"] = carouselLimit
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), data, "获取轮播图列表成功"))
	}
}

func (controller carousel) Create(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	ctx.ShouldBind(&requestData)
	requestFileId := requestData["file_id"]
	fileId := requestFileId.(float64)
	requestTitle := requestData["title"]
	title := requestTitle.(string)
	requestDescription := requestData["description"]
	description := requestDescription.(string)
	requestLink := requestData["link"]
	link := requestLink.(string)
	requestOrder := requestData["order"]
	order := requestOrder.(float64)
	carousel, err := service.CarouselService.Create(int(fileId), title, description, link, int(order))
	if err != nil {
		ctx.JSON(http.StatusOK, controller.ResponseFail(nil, err.Error()))
		return
	}
	carousel.Url = image.FullUrl(carousel.Url)
	carousel.Thumb = image.FullUrl(carousel.Thumb)
	ctx.JSON(http.StatusOK, controller.ResponseOk(carousel, "success"))
}

func (controller carousel) Delete(ctx *gin.Context) {
	requestData := make(map[string]int)
	ctx.ShouldBind(&requestData)
	err := service.CarouselService.Delete(requestData["id"])
	if err != nil {
		ctx.JSON(http.StatusOK, controller.ResponseFail(nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, controller.ResponseOk(nil, "success"))
}

func (controller carousel) Preview(ctx *gin.Context) {
	requestData := make(map[string]int)
	ctx.ShouldBind(&requestData)
	preview, err := service.CarouselService.Preview(requestData["file_id"])
	if err != nil {
		ctx.JSON(http.StatusOK, controller.ResponseFail(nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, controller.ResponseOk(preview, "success"))
}
