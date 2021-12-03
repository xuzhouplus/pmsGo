package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pmsGo/lib/config"
	"pmsGo/lib/controller"
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/middleware/auth"
	"pmsGo/lib/oauth"
	"pmsGo/lib/oauth/gateway"
	"pmsGo/lib/security/json"
	"pmsGo/model"
	"pmsGo/service"
)

const AdminAuthorizeSessionKey = "admin_authorize_session"

type admin struct {
	controller.AppController
}

var Admin = &admin{}

// Verbs 配置方法请求方式
func (ctl admin) Verbs() map[string][]string {
	verbs := make(map[string][]string)
	verbs["login"] = []string{controller.Post}
	verbs["auth"] = []string{controller.Get}
	verbs["logout"] = []string{controller.Post}
	verbs["profile"] = []string{controller.Post, controller.Get}
	verbs["connects"] = []string{controller.Get}
	return verbs
}

// Authenticator 配置方法登录限制
func (ctl admin) Authenticator() controller.Authenticator {
	authenticator := controller.Authenticator{
		Excepts:   []string{},
		Optionals: []string{"login", "auth", "logout"},
	}
	return authenticator
}

// Login 账号登录
func (ctl admin) Login(ctx *gin.Context) {
	requestData := make(map[string]string)
	err := ctx.ShouldBind(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeFail, nil, err.Error()))
		return
	}
	loginAdmin, err := service.AdminService.Login(requestData["account"], requestData["password"])
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeFail, nil, err.Error()))
	} else {
		session := sessions.Default(ctx)
		data, _ := json.Encode(loginAdmin)
		session.Set(auth.SessionLoginAdminKey, data)
		err := session.Save()
		if err != nil {
			ctx.JSON(http.StatusOK, ctl.Response(controller.CodeFail, nil, err.Error()))
			return
		}
		returnAttr := make(map[string]string)
		returnAttr["uuid"] = loginAdmin.Uuid
		returnAttr["type"] = loginAdmin.Type
		returnAttr["avatar"] = loginAdmin.Avatar
		returnAttr["account"] = loginAdmin.Account
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnAttr, "登录成功"))
	}
}

// Auth 获取登录账号信息
func (ctl admin) Auth(ctx *gin.Context) {
	loginAdmin := make(map[string]interface{})
	loginData, _ := ctx.Get(auth.ContextLoginAdminKey)
	if loginData == nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	loginAdmin = loginData.(map[string]interface{})
	if loginAdmin != nil {
		returnAttr := make(map[string]interface{})
		returnAttr["uuid"] = loginAdmin["uuid"]
		returnAttr["type"] = loginAdmin["type"]
		returnAttr["avatar"] = loginAdmin["avatar"]
		returnAttr["account"] = loginAdmin["account"]
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnAttr, "获取成功"))
	} else {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
	}
}

// Logout 登录退出
func (ctl admin) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	err := session.Save()
	if err != nil {
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, nil, "退出成功"))
}

// Profile 登录账号信息编辑
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
		err := ctx.ShouldBind(&requestData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
			return
		}
		if requestData["avatar"] != nil {
			instance, err := fileLib.Base64Upload(requestData["avatar"].(string), "/avatar")
			if err != nil {
				ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
				return
			}
			requestData["avatar"] = string(instance.Url())
			if account.Avatar != "" {
				err := fileLib.Remove(string(fileLib.UrlToPath(fileLib.Url(account.Avatar))))
				if err != nil {
					ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
					return
				}
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

// Connects 获取登录账号绑定第三方信息
func (ctl admin) Connects(ctx *gin.Context) {
	loginAdmin := make(map[string]interface{})
	loginData, _ := ctx.Get(auth.ContextLoginAdminKey)
	if loginData == nil {
		ctx.JSON(http.StatusUnauthorized, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	loginSettings := service.SettingService.GetLoginSettings()
	if len(loginSettings) == 0 {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "获取失败"))
		return
	}
	loginAdmin = loginData.(map[string]interface{})
	connects, err := service.AdminService.GetBoundConnects(loginAdmin["account"].(string), loginSettings)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, map[string]interface{}{
			"types":    loginSettings,
			"connects": nil,
		}, "获取失败"))
		return
	}
	returnData := make(map[string]model.Connect)
	if len(connects) > 0 {
		for _, connect := range connects {
			returnData[connect.Type] = connect
		}
	}
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, map[string]interface{}{
		"types":    loginSettings,
		"connects": returnData,
	}, "获取成功"))
}

// Url 获取授权地址
func (ctl admin) Url(ctx *gin.Context) {
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
	session.Set(AdminAuthorizeSessionKey, encode)
	err = session.Save()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, authorizeUrl, "获取成功"))
}

// Callback oauth授权回调
func (ctl admin) Callback(ctx *gin.Context) {
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
	sessionData := session.Get(AdminAuthorizeSessionKey)
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
	if authorizeData["gateway"].(string) == gateway.TwitterGatewayType {
		if callbackData["oauth_token"] != authorizeData["state"].(string) {
			log.Println("授权数据oauth_token无效")
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
			return
		}
	} else {
		if callbackData["state"] != authorizeData["state"] {
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
	token, err := oauthGateway.AccessToken(callbackData, authorizeData["redirect"].(string))
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
		session := sessions.Default(ctx)
		data, _ := json.Encode(admin)
		session.Set(auth.SessionLoginAdminKey, data)
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
	case "bind":
		contextLoginAdmin, exists := ctx.Get(auth.ContextLoginAdminKey)
		if (!exists) || contextLoginAdmin == nil {
			ctx.JSON(http.StatusUnauthorized, ctl.Response(controller.CodeOk, nil, "需要登录"))
			return
		}
		bind, err := service.AdminService.Bind(int(authorizeData["admin"].(float64)), user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, bind, "获取成功"))
	default:
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, nil, "请求错误"))
	}
}
