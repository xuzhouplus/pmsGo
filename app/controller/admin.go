package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/controller"
	"pmsGo/lib/helper/image"
	"pmsGo/lib/oauth"
	"pmsGo/lib/oauth/gateway"
	"pmsGo/lib/security/json"
	"pmsGo/middleware/auth"
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
		data, _ := json.Encode(loginAdmin)
		session.Set(auth.SessionLoginAdminKey, data)
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
	loginData, _ := ctx.Get(auth.ContextLoginAdminKey)
	if loginData == nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, returnAttr, "获取失败"))
		return
	}
	loginAdmin = loginData.(map[string]interface{})
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
	loginData, _ := ctx.Get(auth.ContextLoginAdminKey)
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
			if account.Avatar != "" {
				image.Remove(string(image.UrlToPath(image.Url(account.Avatar))))
			}
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
	loginData, _ := ctx.Get(auth.ContextLoginAdminKey)
	if loginData == nil {
		ctx.JSON(http.StatusUnauthorized, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	loginAdmin = loginData.(map[string]interface{})
	connects, err := service.AdminService.GetBoundConnects(loginAdmin["account"].(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	returnData := make(map[string]model.Connect)
	if len(connects) > 0 {
		for _, connect := range connects {
			returnData[connect.Type] = connect
		}
	}
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnData, "获取成功"))
}

func (ctl admin) AuthorizeUrl(ctx *gin.Context) {
	gatewayType := ctx.Query("type")
	oauthGateway, err := oauth.NewOauth(gatewayType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	redirect := config.Config.Web.Host + "/profile/authorize/" + gatewayType
	authorizeUrl, state, err := oauthGateway.AuthorizeUrl(ctx.Query("scope"), redirect)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	session := sessions.Default(ctx)
	loginData, _ := ctx.Get(auth.ContextLoginAdminKey)
	adminId := 0
	if loginData != nil {
		loginAdmin := loginData.(map[string]interface{})
		if loginAdmin["id"] != nil {
			adminId = int(loginAdmin["id"].(float64))
		}
	}
	authData := make(map[string]interface{})
	authData["action"] = ctx.Query("action")
	authData["gateway"] = gatewayType
	authData["redirect"] = redirect
	authData["admin"] = adminId
	authData["state"] = state
	encode, err := json.Encode(authData)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	session.Set("authorize", encode)
	err = session.Save()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, authorizeUrl, "获取成功"))
}

func (ctl admin) AuthorizeUser(ctx *gin.Context) {
	requestData := make(map[string]string)
	err := ctx.ShouldBindQuery(&requestData)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
		return
	}
	session := sessions.Default(ctx)
	sessionData := session.Get("authorize")
	if sessionData == nil {
		log.Println("授权数据为空")
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
		return
	}
	authorizeData := make(map[string]interface{})
	err = json.Decode(sessionData.(string), &authorizeData)
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
	if authorizeData["gateway"].(string) == gateway.TwitterGatewayType {
		if requestData["oauth_token"] != authorizeData["state"].(string) {
			log.Println("授权数据oauth_token无效")
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
			return
		}
	} else {
		if requestData["state"] != authorizeData["state"].(string) {
			log.Println("授权数据state无效")
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
			return
		}
	}
	oauthGateway, err := oauth.NewOauth(authorizeData["gateway"].(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	token, err := oauthGateway.AccessToken(requestData, authorizeData["redirect"].(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	user, err := oauthGateway.User(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	switch authorizeData["action"].(string) {
	case "login":
		admin, err := service.AdminService.Auth(user)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "登录失败"))
			return
		}
		if admin == nil {
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "登录失败"))
			return
		}
		loginAdmin := make(map[string]interface{})
		session := sessions.Default(ctx)
		data, _ := json.Encode(loginAdmin)
		session.Set(auth.SessionLoginAdminKey, data)
		session.Save()
		returnAttr := make(map[string]string)
		returnAttr["uuid"] = admin.Uuid
		returnAttr["type"] = admin.Type
		returnAttr["avatar"] = admin.Avatar
		returnAttr["account"] = admin.Account
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnAttr, "登录成功"))
	case "bind":
		contextLoginAdmin, exists := ctx.Get(auth.ContextLoginAdminKey)
		if (!exists) || contextLoginAdmin == nil {
			ctx.JSON(http.StatusUnauthorized, ctl.Response(controller.CodeOk, nil, "需要登录"))
			return
		}
		bind, err := service.AdminService.Bind(int(authorizeData["admin"].(float64)), user)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "绑定失败"))
			return
		}
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, bind, "获取成功"))
	default:
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
	}
}
