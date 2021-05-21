package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/controller"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
	"pmsGo/lib/security/rsa"
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
	if ctx.Request.Method == "GET" {
		result := service.SettingService.Find(model.SiteSettingModel.Keys(), "key")
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}

func (setting setting) Carousel(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		list := service.SettingService.Find(model.CarouselSettingModel.Keys(), "key")
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, list, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}

func (setting setting) Alipay(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.AlipaySettingModel.Keys())
		if result[model.SettingKeyAlipayAppPrimaryKey] !="" {
			decrypt, err := encrypt.Decrypt(base64.Decode(result[model.SettingKeyAlipayAppPrimaryKey]), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				return
			}
			result[model.SettingKeyAlipayAppPrimaryKey] = string(decrypt)
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyAlipayAppPrimaryKey] != nil {
			primaryKey, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyAlipayAppPrimaryKey])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			primaryKeyByte, err := encrypt.Encrypt([]byte(primaryKey), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyAlipayAppPrimaryKey] = base64.Encode(primaryKeyByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}

func (setting setting) Baidu(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.BaiduSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyBaiduSecretKey] != nil {
			primaryKey, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyBaiduSecretKey])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			primaryKeyByte, err := encrypt.Encrypt([]byte(primaryKey), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyBaiduSecretKey] = base64.Encode(primaryKeyByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}

func (setting setting) Facebook(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.FacebookSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyFacebookAppSecret] != nil {
			appSecret, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyFacebookAppSecret])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			appSecretByte, err := encrypt.Encrypt([]byte(appSecret), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyFacebookAppSecret] = base64.Encode(appSecretByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}
func (setting setting) Github(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.GithubSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyGithubAppSecret] != nil {
			appSecret, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyGithubAppSecret])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			appSecretByte, err := encrypt.Encrypt([]byte(appSecret), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyGithubAppSecret] = base64.Encode(appSecretByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}
func (setting setting) Google(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.GoogleSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyGoogleAppSecret] != nil {
			appSecret, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyGoogleAppSecret])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			appSecretByte, err := encrypt.Encrypt([]byte(appSecret), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyGoogleAppSecret] = base64.Encode(appSecretByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}
func (setting setting) Line(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.LineSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyLineAppSecret] != nil {
			appSecret, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyLineAppSecret])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			appSecretByte, err := encrypt.Encrypt([]byte(appSecret), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyLineAppSecret] = base64.Encode(appSecretByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}
func (setting setting) Qq(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.QqSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyQqAppSecret] != nil {
			appSecret, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyQqAppSecret])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			appSecretByte, err := encrypt.Encrypt([]byte(appSecret), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyQqAppSecret] = base64.Encode(appSecretByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}
func (setting setting) Twitter(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.TwitterSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyTwitterAppSecret] != nil {
			appSecret, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyTwitterAppSecret])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			appSecretByte, err := encrypt.Encrypt([]byte(appSecret), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyTwitterAppSecret] = base64.Encode(appSecretByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}
func (setting setting) Wechat(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.WechatSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyWechatAppSecret] != nil {
			appSecret, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyWechatAppSecret])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			appSecretByte, err := encrypt.Encrypt([]byte(appSecret), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyWechatAppSecret] = base64.Encode(appSecretByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
}
func (setting setting) Weibo(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		result := service.SettingService.GetSettings(model.WeiboSettingModel.Keys())
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, result, "获取成功"))
	} else {
		var keyPairs = make(map[string]interface{})
		ctx.ShouldBind(&keyPairs)
		if keyPairs[model.SettingKeyWeiboAppSecret] != nil {
			appSecret, err := rsa.DecryptByPrivateKey(keyPairs[model.SettingKeyWeiboAppSecret])
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			appSecretByte, err := encrypt.Encrypt([]byte(appSecret), []byte(config.Config.Web.Security["salt"]))
			if err != nil {
				ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
				return
			}
			keyPairs[model.SettingKeyWeiboAppSecret] = base64.Encode(appSecretByte)
		}
		err := service.SettingService.Save(keyPairs)
		if err != nil {
			ctx.JSON(http.StatusOK, setting.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, setting.Response(controller.CodeOk, nil, "保存成功"))
	}
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
