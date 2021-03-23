package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pmsGo/lib/helper"
	"strings"
)

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestURI := ctx.Request.RequestURI
		pathSplit := strings.Split(requestURI, "/")
		var requestController string
		var requestAction string
		pathLen := len(pathSplit)
		if pathLen < 2 {
			requestController = "index"
			requestAction = "index"
		} else if pathLen == 2 {
			requestController = pathSplit[1]
			requestAction = "index"
		} else {
			requestController = pathSplit[1]
			requestAction = pathSplit[2]
		}
		controllerId := helper.FirstToUpper(requestController)
		actionId := helper.FirstToUpper(requestAction)
		fmt.Printf("controller:%v action:%v \n", controllerId, actionId)
		fmt.Println(ctx.Handler())
	}
}
