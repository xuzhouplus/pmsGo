package auth

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pmsGo/lib/config"
	"pmsGo/lib/helper"
	"strings"
)

var settings map[interface{}]interface{}

func init() {
	settings = config.Config.Web["auth"].(map[interface{}]interface{})
	log.Printf("auth: %v \n", settings)
}

func getRequest(ctx *gin.Context) (string, string) {
	var controller string
	var action string
	uri := ctx.Request.RequestURI
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
	authSet := settings[controller]
	if authSet == nil {
		return nil
	}
	authSetMap := authSet.(map[interface{}]interface{})
	except := authSetMap["except"]
	if except != nil {
		exceptSet := except.([]interface{})
		_, result := helper.IsInSlice(exceptSet, action)
		if result {
			return "except"
		}
	}
	optional := authSetMap["optional"]
	if optional != nil {
		optionalSet := optional.([]interface{})
		_, result := helper.IsInSlice(optionalSet, action)
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
			sessionAdmin := session.Get("login_admin")
			if sessionAdmin == nil && authType == nil {
				ctx.JSON(http.StatusUnauthorized, nil)
				ctx.Abort()
				return
			}
			if sessionAdmin != nil {
				loginAdmin := make(map[string]interface{})
				err := json.Unmarshal(sessionAdmin.([]byte), &loginAdmin)
				if err!=nil {
					log.Printf("解析session数据失败,%e", err)
				}else {
					ctx.Set("loginAdmin", loginAdmin)
				}
			}
		}
		ctx.Next()
	}
}
