package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/lib/controller"
	fileHelper "pmsGo/lib/helper/file"
)

type file struct {
	controller.App
}

var File = &file{}

func (cto file) Upload(ctx *gin.Context) {
	upload:=&fileHelper.Upload{}
 err := upload.Upload(ctx, "file","/file")
	if err != nil {
		ctx.JSON(http.StatusOK,cto.ResponseFail("", err.Error()))
	} else {
		ctx.JSON(http.StatusOK,cto.ResponseOk(upload, "success"))
	}
}
