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
	verbs["view"] = []string{controller.Get}
	verbs["set-caption-style"] = []string{controller.Post}
	return verbs
}

func (ctl carousel) Authenticator() controller.Authenticator {
	authenticator := controller.Authenticator{
		Excepts: []string{"index"},
	}
	return authenticator
}
func (ctl carousel) Index(ctx *gin.Context) {
	match := map[string]interface{}{
		"status": model.CarouselStatusEnabled,
	}
	carousels, err := service.CarouselService.List(0, 0, match, "", nil)
	if err != nil {
		log.Error(err)
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, err.Error()))
	} else {
		result := make([]map[string]interface{}, 0)
		for _, carousel := range carousels {
			result = append(result, map[string]interface{}{
				"uuid":        carousel.Uuid,
				"type":        carousel.Type,
				"title":       carousel.Title,
				"description": carousel.Description,
				"url":         fileLib.FullUrl(carousel.Url),
				"width":       carousel.Width,
				"height":      carousel.Height,
				"link":        carousel.Link,
				"switch_type": carousel.SwitchType,
				"timeout":     carousel.Timeout,
			})
		}
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeOk(), result, "获取轮播图列表成功"))
	}
}

func (ctl carousel) List(ctx *gin.Context) {
	carousels, err := service.CarouselService.List(0, 0, nil, "", nil)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, err.Error()))
	} else {
		data := make(map[string]interface{})
		for index, carousel := range carousels {
			carousels[index].Url = fileLib.FullUrl(carousel.Url)
			carousels[index].Thumb = fileLib.FullUrl(carousel.Thumb)
		}
		data["list"] = carousels
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
	requestTimeout := requestData["timeout"]
	timeout := requestTimeout.(float64)
	carousel, err := service.CarouselService.Create(int(fileId), title, description, link, int(order), switchType, int(timeout))
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
	requestTimeout := requestData["timeout"]
	timeout := requestTimeout.(float64)
	carousel, err := service.CarouselService.Update(int(id), int(fileId), title, description, link, int(order), switchType, int(timeout))
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
	err := ctx.ShouldBind(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
	err = service.CarouselService.Delete(requestData["id"])
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.ResponseOk(nil, "success"))
}

func (ctl carousel) View(ctx *gin.Context) {
	queryId, _ := ctx.GetQuery("id")
	if queryId == "" {
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, "缺少请求参数"))
		return
	}
	carouselId, err := strconv.Atoi(queryId)
	if err != nil {
		return
	}
	carousel, err := service.CarouselService.FindById(carouselId)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, "缺少请求参数"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeOk(), map[string]interface{}{
		"uuid":   carousel.Uuid,
		"url":    fileLib.FullUrl(carousel.Url),
		"thumb":  fileLib.FullUrl(carousel.Thumb),
		"status": carousel.Status,
	}, "获取成功"))
}

func (ctl carousel) SetCaptionStyle(ctx *gin.Context) {
	requestData := &model.CaptionStyle{}
	err := ctx.ShouldBind(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.ResponseFail(nil, err.Error()))
		return
	}
	queryId, _ := ctx.GetQuery("id")
	if queryId == "" {
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, "缺少请求参数"))
		return
	}
	carouselId, _ := strconv.Atoi(queryId)
	carouselField, _ := ctx.GetQuery("field")
	if carouselField == "" {
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, "缺少请求参数"))
		return
	}
	if carouselField == "title" {
		err = service.CarouselService.SetTitleStyle(carouselId, requestData)
	} else {
		err = service.CarouselService.SetDescriptionStyle(carouselId, requestData)
	}
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeFail(), nil, "保存失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.Response(ctl.CodeOk(), nil, "设置成功"))
}
