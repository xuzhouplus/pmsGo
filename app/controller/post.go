package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/controller"
	"pmsGo/lib/image"
	"strconv"
)

type post struct {
	controller.App
}

var Post = &post{}

func (controller post) Index(ctx *gin.Context) {
	page := ctx.Query("page")
	pageNum, _ := strconv.Atoi(page)
	limit := ctx.Query("limit")
	limitNum, _ := strconv.Atoi(limit)
	search := ctx.Query("search")
	order := ctx.QueryMap("order")
	posts, err := service.PostService.List(pageNum, limitNum, []string{"uuid", "cover", "title", "sub_title", "created_at", "updated_at"}, search, 1, order)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), posts, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), posts, "获取稿件列表成功"))
	}
}

func (controller post) List(ctx *gin.Context) {
	page := ctx.Query("page")
	pageNum, _ := strconv.Atoi(page)
	limit := ctx.Query("limit")
	limitNum, _ := strconv.Atoi(limit)
	search := ctx.Query("search")
	enable := ctx.Query("enable")
	enableNum, _ := strconv.Atoi(enable)
	order := ctx.QueryMap("order")
	posts, err := service.PostService.List(pageNum, limitNum, []string{}, search, enableNum, order)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), posts, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), posts, "获取稿件列表成功"))
	}
}

func (controller post) Save(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	ctx.ShouldBind(&requestData)
	save, err := service.PostService.Save(requestData)
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), save, "创建稿件列表成功"))
}
func (controller post) Delete(ctx *gin.Context) {
	requestData := make(map[string]int)
	ctx.ShouldBind(&requestData)
	err := service.PostService.Delete(requestData["id"])
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), nil, "删除稿件列表成功"))
}
func (controller post) ToggleStatus(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	ctx.ShouldBind(&requestData)
	var one *model.Post
	var err error
	switch requestData["id"].(type) {
	case string:
		one, err = service.PostService.FindOneByUuid(requestData["id"].(string))
	case int:
		one, err = service.PostService.FindOneById(requestData["id"].(int))
	case float64:
		one, err = service.PostService.FindOneById(int(requestData["id"].(float64)))
	}
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
		return
	}
	if one == nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, "稿件不存在"))
		return
	}
	err = one.Toggle()
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), one, "修改成功"))
}

func (controller post) Detail(ctx *gin.Context) {
	var one *model.Post
	var err error
	id := ctx.Query("id")
	if len(id) != 32 {
		idNum, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
			return
		}
		one, err = service.PostService.FindOneById(idNum)
	} else {
		one, err = service.PostService.FindOneByUuid(id)
	}
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
		return
	}
	one.Cover = image.FullUrl(one.Cover)
	ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), one, "获取稿件成功"))
}
