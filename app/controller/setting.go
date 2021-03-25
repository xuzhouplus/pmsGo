package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/lib/controller"
)

type setting struct {
	controller.App
}

var Setting = &setting{}

func (setting setting) Index(c *gin.Context) {
	result, err := model.SettingModel.GetPublicSettings()
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
	result := model.SiteSettingModel.Find(model.SiteSettingModel.Keys(), "key")
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}

func (setting setting) Carousel(ctx *gin.Context) {
	list := model.CarouselSettingModel.Find(model.CarouselSettingModel.Keys(), "key")
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, list, "获取成功"))
}

func (setting setting) Alipay(ctx *gin.Context) {
	result := model.AlipaySettingModel.GetSettings(model.AlipaySettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}

func (setting setting) SaveAlipay(ctx *gin.Context) {
	var keyPairs = make(map[string]interface{})
	ctx.ShouldBind(&keyPairs)
	model.CarouselSettingModel.Save(keyPairs)
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
}
func (setting setting) Baidu(ctx *gin.Context) {
	result := model.BaiduSettingModel.GetSettings(model.BaiduSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}

func (setting setting) Facebook(ctx *gin.Context) {
	result := model.FacebookSettingModel.GetSettings(model.FacebookSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Github(ctx *gin.Context) {
	result := model.GithubSettingModel.GetSettings(model.GithubSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Google(ctx *gin.Context) {
	result := model.GoogleSettingModel.GetSettings(model.GoogleSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Line(ctx *gin.Context) {
	result := model.LineSettingModel.GetSettings(model.LineSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Qq(ctx *gin.Context) {
	result := model.QqSettingModel.GetSettings(model.QqSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Twitter(ctx *gin.Context) {
	result := model.TwitterSettingModel.GetSettings(model.TwitterSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Wechat(ctx *gin.Context) {
	result := model.WechatSettingModel.GetSettings(model.WechatSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}
func (setting setting) Weibo(ctx *gin.Context) {
	result := model.WeiboSettingModel.GetSettings(model.WeiboSettingModel.Keys())
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
}

func (setting setting) Save(ctx *gin.Context) {
	var keyPairs = make(map[string]interface{})
	ctx.ShouldBind(&keyPairs)
	model.CarouselSettingModel.Save(keyPairs)
	ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
}
