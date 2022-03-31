package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/lib/controller"
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/log"
	"pmsGo/model"
	"pmsGo/service"
	"strconv"
)

type carousel struct {
	controller.AppController
}

var Carousel = &carousel{}

func (ctl carousel) Verbs() map[string][]string {
	verbs := make(map[string][]string)
	verbs["index"] = []string{controller.Get}
	verbs["list"] = []string{controller.Get}
	verbs["create"] = []string{controller.Post}
	verbs["update"] = []string{controller.Post}
	verbs["delete"] = []string{controller.Post}
	verbs["preview"] = []string{controller.Post}
	return verbs
}

func (ctl carousel) Authenticator() controller.Authenticator {
	authenticator := controller.Authenticator{
		Excepts: []string{"index"},
	}
	return authenticator
}
func (ctl carousel) Index(ctx *gin.Context) {
	fields, _ := ctx.GetQueryArray("fields")
	like, _ := ctx.GetQuery("like")
	order, _ := ctx.GetQueryMap("order")
	result, err := service.CarouselService.List(0, 0, fields, like, order)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeOk(), result, "获取轮播图列表成功"))
	}
}

func (ctl carousel) List(ctx *gin.Context) {
	list, err := service.CarouselService.List(0, 0, nil, nil, nil)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, err.Error()))
	} else {
		data := make(map[string]interface{})
		data["list"] = list
		carouselLimit, _ := strconv.Atoi(service.SettingService.GetSetting(model.SettingKeyCarouselLimit))
		data["limit"] = carouselLimit
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeOk(), data, "获取轮播图列表成功"))
	}
}

func (ctl carousel) Create(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	err := ctx.ShouldBind(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
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
	requestSwitchType := requestData["switch_type"]
	switchType := requestSwitchType.(string)
	carousel, err := service.CarouselService.Create(int(fileId), title, description, link, int(order), switchType)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
	carousel.Url = fileLib.FullUrl(carousel.Url)
	carousel.Thumb = fileLib.FullUrl(carousel.Thumb)
	ctx.JSON(http.StatusOK, ctl.ResponseOk(carousel, "success"))
}

func (ctl carousel) Update(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	err := ctx.ShouldBind(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
	requestId := requestData["id"]
	id := requestId.(float64)
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
	requestSwitchType := requestData["switch_type"]
	switchType := requestSwitchType.(string)
	carousel, err := service.CarouselService.Update(int(id), int(fileId), title, description, link, int(order), switchType)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
	carousel.Url = fileLib.FullUrl(carousel.Url)
	carousel.Thumb = fileLib.FullUrl(carousel.Thumb)
	ctx.JSON(http.StatusOK, ctl.ResponseOk(carousel, "success"))
}

func (ctl carousel) Delete(ctx *gin.Context) {
	requestData := make(map[string]int)
	ctx.ShouldBind(&requestData)
	err := service.CarouselService.Delete(requestData["id"])
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.ResponseOk(nil, "success"))
}

func (ctl carousel) Preview(ctx *gin.Context) {
	requestData := make(map[string]int)
	ctx.ShouldBind(&requestData)
	preview, err := service.CarouselService.Preview(requestData["file_id"])
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.ResponseOk(preview, "success"))
}
