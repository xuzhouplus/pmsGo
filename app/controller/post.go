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
