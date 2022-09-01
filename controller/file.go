package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/lib/controller"
	fileLib "pmsGo/lib/file"
	"pmsGo/service"
	"strconv"
)

type file struct {
	controller.AppController
}

var File = &file{}

func (cto file) Verbs() map[string][]string {
	verbs := make(map[string][]string)
	verbs["index"] = []string{controller.Get}
	verbs["upload"] = []string{controller.Post, controller.Get}
	verbs["delete"] = []string{controller.Post}
	verbs["detail"] = []string{controller.Get}
	verbs["extractFrame"] = []string{controller.Post}
	verbs["getFrame"] = []string{controller.Get}
	return verbs
}

func (cto file) Authenticator() controller.Authenticator {
	authenticator := controller.Authenticator{
		Excepts:   []string{},
		Optionals: []string{"index", "upload", "delete"},
	}
	return authenticator
}

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
	if ctx.Request.Method == controller.Get {
		if fileLib.ChunkCheck(ctx) == true {
			ctx.JSON(http.StatusOK, cto.ResponseOk("", ""))
		}
		ctx.JSON(http.StatusNotFound, cto.ResponseFail("", ""))
		return
	}
	formUpload, err := fileLib.ChunkUpload(ctx, "binary", fileLib.SubDir)
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	if formUpload.Status == fileLib.UploadStatusProcess {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", ""))
		return
	}
	fileModel, err := service.FileService.Upload(formUpload, ctx.PostForm("name"), ctx.PostForm("description"))
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	switch formUpload.FileType {
	case fileLib.TypeVideo:
		err := service.FileService.ProcessVideo(fileModel)
		if err != nil {
			ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
			return
		}
	case fileLib.TypeImage:
		err := service.FileService.ProcessImage(fileModel)
		if err != nil {
			ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
			return
		}
	}
	ctx.JSON(http.StatusOK, cto.ResponseOk(fileModel, "success"))
}

func (cto file) ExtractFrame(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	err := ctx.ShouldBindJSON(&requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	taskId, err := service.FileService.ExtractFrame(requestData["file_id"].(string), requestData["point"].(int64), requestData["width"].(int64), requestData["height"].(int64))
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, cto.ResponseOk(taskId, "发起抽帧成功"))
}

func (cto file) GetFrame(ctx *gin.Context) {
	taskId := ctx.Query("task_id")
	if taskId == "" {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", "缺少任务参数"))
		return
	}
	result := service.FileService.GetFrame(taskId)
	ctx.JSON(http.StatusOK, cto.ResponseOk(result, "发起抽帧成功"))
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

func (cto file) Detail(ctx *gin.Context) {
	id := ctx.Query("uuid")
	one, err := service.FileService.FindByUuid(id)
	if err != nil {
		ctx.JSON(http.StatusOK, cto.ResponseFail("", err.Error()))
		return
	}
	one.Path = fileLib.FullUrl(one.Path)
	one.Thumb = fileLib.FullUrl(one.Thumb)
	one.Preview = fileLib.FullUrl(one.Preview)
	ctx.JSON(http.StatusOK, cto.ResponseOk(one, "success"))
}
