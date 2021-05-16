package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/lib/controller"
	fileHelper "pmsGo/lib/helper/image"
)

type file struct {
	controller.App
}

var File = &file{}

func (cto file) Index(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	err := ctx.ShouldBind(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	list, err := model.FileModel.List(requestData["page"], requestData["limit"], requestData["select"], requestData["type"], requestData["name"])
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, cto.ResponseOk(list, "success"))
}

func (cto file) Upload(ctx *gin.Context) {
	upload := &fileHelper.Upload{}
	err := upload.Upload(ctx, "file", "/image")
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
	} else {
		model.FileModel.Upload(upload, ctx.PostForm("name"), ctx.PostForm("description"))
		ctx.JSON(http.StatusOK, cto.ResponseOk(model.FileModel, "success"))
	}
}

func (cto file) Delete(ctx *gin.Context) {
	requestData := make(map[string]int)
	err := ctx.ShouldBindJSON(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	err = model.FileModel.Delete(requestData["id"])
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, cto.ResponseOk(nil, "success"))
}
