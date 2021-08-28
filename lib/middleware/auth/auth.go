package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/lib/controller"
	"pmsGo/lib/log"
	"pmsGo/lib/security/json"
	"strings"
)

// SessionLoginAdminKey 登录信息session保存key
const SessionLoginAdminKey = "login_admin"

// ContextLoginAdminKey 登录信息上下文保存key
const ContextLoginAdminKey = "loginAdmin"

var Authenticators = make(map[string]*controller.Resolve)

//获取请求controller和action
func getRequest(ctx *gin.Context) (string, string) {
	var controllerName string
	var actionName string
	uri := ctx.Request.URL.Path
	splits := strings.Split(uri, "/")
	leng := len(splits)
	if leng < 2 {
		controllerName = "index"
		actionName = "index"
	} else if leng == 2 {
		controllerName = splits[1]
		actionName = "index"
	} else {
		controllerName = splits[1]
		actionName = splits[2]
	}
	return controllerName, actionName
}

//获取授权类型
func getAuthType(controllerName string, actionName string) string {
	controllerAuthenticator := Authenticators[controllerName]
	if controllerAuthenticator == nil {
		return controller.Forbidden
	}
	actionAuthenticator := controllerAuthenticator.GetAction(actionName)
	if actionAuthenticator == nil {
		return controller.Forbidden
	}
	return actionAuthenticator.Authenticator
}

// Add 添加验证数据
func Add(resolve *controller.Resolve) {
	Authenticators[resolve.GetControllerName()] = resolve
}

// Auth Register 注册登录判断中间件
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取请求controller和action
		controllerName, actionName := getRequest(ctx)
		//获取认证配置
		authType := getAuthType(controllerName, actionName)
		log.Debugf("controller: %s ,action: %s ,auth: %s\n", controllerName, actionName, authType)
		//如果不是不需要登录
		if authType != controller.Except {
			//获取登录用户信息
			session := sessions.Default(ctx)
			sessionAdmin := session.Get(SessionLoginAdminKey)
			if sessionAdmin == nil {
				//用户登录信息为空且认证类型不为可不登录，直接返回401
				if authType == controller.Forbidden {
					ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
						"code":    0,
						"data":    nil,
						"message": "Unauthorized",
					})
					ctx.Abort()
					return
				}
			} else {
				//有登录信息，把登录信息写入上下文环境中
				loginAdmin := make(map[string]interface{})
				err := json.Decode(sessionAdmin.(string), &loginAdmin)
				if err != nil {
					log.Errorf("解析session数据失败,%e", err)
				} else {
					if loginAdmin != nil {
						ctx.Set(ContextLoginAdminKey, loginAdmin)
					}
				}
			}
		}
		ctx.Next()
	}
}
