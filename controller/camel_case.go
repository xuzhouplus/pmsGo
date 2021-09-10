package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/lib/controller"
)

type camelCase struct {
	controller.AppController
}

var CamelCase = &camelCase{}

func (receiver camelCase) Index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, receiver.ResponseOk(nil, ""))
}
