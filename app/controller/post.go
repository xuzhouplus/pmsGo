package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/lib/controller"
)

type post struct {
	controller.App
}

var Post = &post{}

func (controller post) Index(c *gin.Context) {
	requestData := make(map[string]interface{})
	c.ShouldBind(&requestData)
	posts, error := model.PostModel.List(requestData["page"], requestData["limit"], []string{"uuid", "cover", "title", "sub_title", "created_at", "updated_at"}, requestData["search"], 1, requestData["order"])
	if error != nil {
		c.JSON(http.StatusOK, controller.Response(controller.CodeFail(), posts, error.Error()))
	} else {
		c.JSON(http.StatusOK, controller.Response(controller.CodeOk(), posts, "获取稿件列表成功"))
	}
}

func (controller post) List(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	ctx.ShouldBind(&requestData)
	posts, err := model.PostModel.List(requestData["page"], requestData["limit"], nil, requestData["search"], 1, requestData["order"])
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), posts, err.Error()))
	} else {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), posts, "获取稿件列表成功"))
	}
}

func (controller post) Save(ctx *gin.Context) {

}
func (controller post) Delete(ctx *gin.Context) {

}
func (controller post) ToggleStatus(ctx *gin.Context) {
	requestData := make(map[string]interface{})
	ctx.ShouldBind(&requestData)
	var one *model.Post
	var err error
	switch requestData["id"].(type) {
	case string:
		one, err = model.PostModel.FindOneByUuid(requestData["id"].(string))
	case int:
		one, err = model.PostModel.FindOneById(requestData["id"].(int))
	}
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
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
	requestData := make(map[string]interface{})
	ctx.ShouldBind(&requestData)
	var one *model.Post
	var err error
	switch requestData["id"].(type) {
	case string:
		one, err = model.PostModel.FindOneByUuid(requestData["id"].(string))
	case int:
		one, err = model.PostModel.FindOneById(requestData["id"].(int))
	}
	if err != nil {
		ctx.JSON(http.StatusOK, controller.Response(controller.CodeFail(), nil, err.Error()))
	}
	ctx.JSON(http.StatusOK, controller.Response(controller.CodeOk(), one, "获取稿件成功"))
}
