package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/controller"
	"pmsGo/lib/helper/image"
	"pmsGo/lib/oauth"
	"pmsGo/lib/oauth/gateway"
	"pmsGo/lib/security/random"
)

type admin struct {
	controller.App
}

var Admin = &admin{}

func (ctl admin) Login(ctx *gin.Context) {
	requestData := make(map[string]string)
	ctx.ShouldBind(&requestData)
	loginAdmin, err := service.AdminService.Login(requestData["account"], requestData["password"])
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeFail, nil, err.Error()))
	} else {
		session := sessions.Default(ctx)
		data, _ := json.Marshal(loginAdmin)
		session.Set("login_admin", data)
		session.Save()
		returnAttr := make(map[string]string)
		returnAttr["uuid"] = loginAdmin.Uuid
		returnAttr["type"] = loginAdmin.Type
		returnAttr["avatar"] = loginAdmin.Avatar
		returnAttr["account"] = loginAdmin.Account
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnAttr, "登录成功"))
	}
}

func (ctl admin) Auth(ctx *gin.Context) {
	returnAttr := make(map[string]interface{})
	loginAdmin := make(map[string]interface{})
	session := sessions.Default(ctx)
	loginData := session.Get("login_admin")
	if loginData == nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, returnAttr, "获取失败"))
		return
	}
	err := json.Unmarshal(loginData.([]byte), &loginAdmin)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, returnAttr, "获取失败"))
	}
	if loginAdmin != nil {
		returnAttr["uuid"] = loginAdmin["uuid"]
		returnAttr["type"] = loginAdmin["type"]
		returnAttr["avatar"] = loginAdmin["avatar"]
		returnAttr["account"] = loginAdmin["account"]
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnAttr, "获取成功"))
	} else {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, returnAttr, "获取失败"))
	}
}

func (ctl admin) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, nil, "退出成功"))
}

func (ctl admin) Profile(ctx *gin.Context) {
	loginAdmin := make(map[string]interface{})
	loginData, _ := ctx.Get("loginAdmin")
	if loginData == nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	loginAdmin = loginData.(map[string]interface{})
	if loginAdmin == nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	account, err := service.AdminService.FindOneByAccount(loginAdmin["account"].(string), 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	if ctx.Request.Method == "GET" {
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, account, "获取成功"))
	} else {
		requestData := make(map[string]interface{})
		ctx.ShouldBind(&requestData)
		if requestData["avatar"] != nil {
			instance, err := image.Base64Upload(requestData["avatar"].(string), "/avatar")
			if err != nil {
				ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
				return
			}
			requestData["avatar"] = string(instance.Url())
		}
		admin, err := service.AdminService.Update(requestData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "保存失败"))
			return
		}
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, admin, "保存成功"))
	}
}

func (ctl admin) Connects(ctx *gin.Context) {
	loginAdmin := make(map[string]interface{})
	session := sessions.Default(ctx)
	loginData := session.Get("login_admin")
	if loginData == nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	err := json.Unmarshal(loginData.([]byte), &loginAdmin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	connects, err := service.AdminService.GetBoundConnects(loginAdmin["account"].(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, connects, "获取成功"))
}

func (ctl admin) AuthorizeUrl(ctx *gin.Context) {
	gatewayType := gateway.GitHubGatewayType
	oauthGateway, err := oauth.NewOauth(gatewayType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	redirect := config.Config.Web.Host + "/profile/authorize/" + gatewayType
	state := random.Uuid(false)
	authorizeUrl := oauthGateway.AuthorizeUrl("", redirect, state)
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, authorizeUrl, "获取成功"))
}

func (ctl admin) AuthorizeUser(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	gatewayType := ctx.Query("gateway")
	oauthGateway, err := oauth.NewOauth(gatewayType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	redirect := config.Config.Web.Host + "/profile/authorize/" + gatewayType
	token, err := oauthGateway.AccessToken(code, redirect, state)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	user, err := oauthGateway.User(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, user, "获取成功"))
}
