package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pmsGo/lib/cache"
	"pmsGo/lib/config"
	"pmsGo/lib/controller"
	"pmsGo/lib/middleware/auth"
	nasLib "pmsGo/lib/nas"
	"pmsGo/lib/security/json"
	"pmsGo/lib/security/random"
	"pmsGo/service"
)

const NasAuthorizeSessionKey = "nas_authorize_session"
const NasAuthorizeCacheKey = "nas_authorize_cache"

type nas struct {
	controller.AppController
}

var Nas = &nas{}

// Verbs 配置方法请求方式
func (ctl nas) Verbs() map[string][]string {
	verbs := make(map[string][]string)
	verbs["url"] = []string{controller.Get}
	verbs["callback"] = []string{controller.Get}
	return verbs
}

// Authenticator 配置方法登录限制
func (ctl nas) Authenticator() controller.Authenticator {
	authenticator := controller.Authenticator{
		Excepts:   []string{},
		Optionals: []string{"url", "callback"},
	}
	return authenticator
}

// Authorize 获取授权地址
func (ctl nas) Authorize(ctx *gin.Context) {
	gatewayType := ctx.Query("type")
	nasGateway, err := nasLib.NewNas(gatewayType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	redirect := config.Config.Web.Host + "/nas/authorize/" + gatewayType
	state := random.Uuid(false)
	authorizeUrl, err := nasGateway.Gateway.Authorize(redirect, state)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	session := sessions.Default(ctx)
	authData := make(map[string]interface{})
	authData["gateway"] = gatewayType
	authData["redirect"] = redirect
	authData["state"] = state
	encode, err := json.Encode(authData)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	session.Set(NasAuthorizeSessionKey, encode)
	err = session.Save()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, authorizeUrl, "获取成功"))
}

// Callback oauth授权回调
func (ctl nas) Callback(ctx *gin.Context) {
	requestData := ctx.Request.URL.Query()
	if len(requestData) == 0 {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
		return
	}
	callbackData := make(map[string]string)
	for field, strings := range requestData {
		callbackData[field] = strings[0]
	}
	session := sessions.Default(ctx)
	sessionData := session.Get(NasAuthorizeSessionKey)
	if sessionData == nil {
		log.Println("授权数据为空")
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
		return
	}
	authorizeData := make(map[string]interface{})
	err := json.Decode(sessionData.(string), &authorizeData)
	if err != nil {
		log.Println("授权数据解析失败")
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
		return
	}
	if authorizeData["state"] == nil {
		log.Println("授权数据state为空")
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
		return
	}
	if authorizeData["gateway"] == nil {
		log.Println("授权数据gateway无效")
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
		return
	}
	if callbackData["state"] != authorizeData["state"] {
		log.Println("授权数据state无效")
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
		return
	}
	nasGateway, err := nasLib.NewNas(authorizeData["gateway"].(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	accessToken, err := nasGateway.Gateway.AccessToken(callbackData, authorizeData["redirect"].(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	authorizeUser,err:=nasGateway.Gateway.User(accessToken.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	data, _ := json.Encode(authorizeUser)
	cache.Set(NasAuthorizeCacheKey+nasGateway.Type,data,0)
	session.Delete(NasAuthorizeSessionKey)
	err = session.Save()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	returnAttr := make(map[string]string)
	returnAttr["uuid"] = admin.Uuid
	returnAttr["type"] = admin.Type
	returnAttr["avatar"] = admin.Avatar
	returnAttr["account"] = admin.Account
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnAttr, "登录成功"))
}
