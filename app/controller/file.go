package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/service"
	"pmsGo/lib/controller"
	"pmsGo/lib/image"
	"strconv"
)

type file struct {
	controller.App
}

var File = &file{}

func (cto file) Index(ctx *gin.Context) {
	page := ctx.Query("page")
	pageNum, _ := strconv.Atoi(page)
	limit := ctx.Query("limit")
	limitNum, _ := strconv.Atoi(limit)
	fields := ctx.QueryArray("select[]")
	fileType := ctx.Query("type")
	name := ctx.Query("name")
	list, err := service.FileService.List(pageNum, limitNum, fields, fileType, name)
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, cto.ResponseOk(list, "success"))
}

func (cto file) Upload(ctx *gin.Context) {
	upload, err := image.Upload(ctx, "file", "/image")
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
	} else {
		fileModel, err := service.FileService.Upload(upload, ctx.PostForm("name"), ctx.PostForm("description"))
		if err != nil {
			ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
			return
		}
		fileModel.Path = image.FullUrl(fileModel.Path)
		fileModel.Thumb = image.FullUrl(fileModel.Thumb)
		fileModel.Preview = image.FullUrl(fileModel.Preview)
		ctx.JSON(http.StatusOK, cto.ResponseOk(fileModel, "success"))
	}
}

func (cto file) Delete(ctx *gin.Context) {
	requestData := make(map[string]int)
	err := ctx.ShouldBindJSON(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	err = service.FileService.Delete(requestData["id"])
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, cto.ResponseOk(nil, "success"))
}
