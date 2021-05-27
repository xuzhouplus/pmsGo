package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pmsGo/lib/config"
	"pmsGo/lib/helper"
	"pmsGo/lib/security/json"
	"strings"
)

var settings map[string]config.Auth

const SessionLoginAdminKey = "login_admin"
const ContextLoginAdminKey = "loginAdmin"

func init() {
	settings = config.Config.Web.Auth
	log.Printf("auth: %v \n", settings)
}

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
func getAuthType(controller string, action string) interface{} {
	authSetMap := settings[controller]
	except := authSetMap.Except
	if except != nil {
		_, result := helper.IsInSlice(except, action)
		if result {
			return "except"
		}
	}
	optional := authSetMap.Optional
	if optional != nil {
		_, result := helper.IsInSlice(optional, action)
		if result {
			return "optional"
		}
	}
	return nil
}
func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		controller, action := getRequest(ctx)
		authType := getAuthType(controller, action)
		log.Printf("controller: %v ,action: %v ,auth: %v\n", controller, action, authType)
		if authType != "except" {
			session := sessions.Default(ctx)
			sessionAdmin := session.Get(SessionLoginAdminKey)
			if sessionAdmin == nil {
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
				loginAdmin := make(map[string]interface{})
				err := json.Decode(sessionAdmin.(string), &loginAdmin)
				if err != nil {
					log.Printf("解析session数据失败,%e", err)
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
