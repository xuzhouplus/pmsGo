package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/controller"
)

type setting struct {
	controller.App
}

var Setting = &setting{}

func (setting setting) Index(c *gin.Context) {
	result, err := service.SettingService.GetPublicSettings()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	returnData := make(map[string]interface{})
	returnData["code"] = 1
	returnData["data"] = result
	c.JSON(http.StatusOK, returnData)
}

func (setting setting) Site(ctx *gin.Context) {
	result := service.SettingService.Find(model.SiteSettingModel.Keys(), "key")
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}

func (setting setting) Carousel(ctx *gin.Context) {
	list := service.SettingService.Find(model.CarouselSettingModel.Keys(), "key")
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, list, "获取成功"))
}

func (setting setting) Alipay(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.AlipaySettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}

func (setting setting) SaveAlipay(ctx *gin.Context) {
	var keyPairs = make(map[string]interface{})
	ctx.ShouldBind(&keyPairs)
	service.SettingService.Save(keyPairs)
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
}
func (setting setting) Baidu(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.BaiduSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}

func (setting setting) Facebook(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.FacebookSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Github(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.GithubSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Google(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.GoogleSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Line(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.LineSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Qq(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.QqSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Twitter(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.TwitterSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Wechat(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.WechatSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Weibo(ctx *gin.Context) {
	result := service.SettingService.GetSettings(model.WeiboSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}

func (setting setting) Save(ctx *gin.Context) {
	var keyPairs = make(map[string]interface{})
	ctx.ShouldBind(&keyPairs)
	err := service.SettingService.Save(keyPairs)
	if err != nil {
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
}
