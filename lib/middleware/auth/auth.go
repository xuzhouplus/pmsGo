package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/lib/config"
	"pmsGo/lib/helper"
	"pmsGo/lib/log"
	"pmsGo/lib/security/json"
	"strings"
)

// SessionLoginAdminKey 登录信息session保存key
const SessionLoginAdminKey = "login_admin"

// ContextLoginAdminKey 登录信息上下文保存key
const ContextLoginAdminKey = "loginAdmin"

// TypeExcept 不需要登录
const TypeExcept = "except"

// TypeOptional 可以不登录
const TypeOptional = "optional"

//获取请求controller和action
func getRequest(ctx *gin.Context) (string, string) {
	var controller string
	var action string
	uri := ctx.Request.URL.Path
	splits := strings.Split(uri, "/")
	leng := len(splits)
	if leng < 2 {
		controller = "index"
		action = "index"
	} else if leng == 2 {
		controller = splits[1]
		action = "index"
	} else {
		controller = splits[1]
		action = splits[2]
	}
	return controller, action
}

//获取请求controller和action配置的授权类型
func getAuthType(controller string, action string) interface{} {
	authSetMap := config.Config.Web.Auth[controller]
	except := authSetMap.Except
	if except != nil {
		_, result := helper.IsInSlice(except, action)
		if result {
			return TypeExcept
		}
	}
	optional := authSetMap.Optional
	if optional != nil {
		_, result := helper.IsInSlice(optional, action)
		if result {
			return TypeOptional
		}
	}
	return nil
}

// Register 注册登录判断中间件
func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取请求controller和action
		controller, action := getRequest(ctx)
		//获取认证配置
		authType := getAuthType(controller, action)
		log.Debugf("controller: %s ,action: %s ,auth: %s\n", controller, action, authType)
		//如果不是不需要登录
		if authType != TypeExcept {
			//获取登录用户信息
			session := sessions.Default(ctx)
			sessionAdmin := session.Get(SessionLoginAdminKey)
			if sessionAdmin == nil {
				//用户登录信息为空且认证类型不为可不登录，直接返回401
				if authType == nil {
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
