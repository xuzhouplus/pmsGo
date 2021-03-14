package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
)

type carousel struct {
	app
}

var Carousel = &carousel{}

func (controller carousel) Index(c *gin.Context) {
	requestData := make(map[string]interface{})
	c.ShouldBind(&requestData)
	result, err := model.Carousel.List(0, 0, requestData["fields"], requestData["like"], requestData["order"])
	if err != nil {
		c.JSON(http.StatusOK, controller.response(CodeOk, result, err.Error()))
	} else {
		c.JSON(http.StatusOK, controller.response(CodeOk, result, "获取轮播图列表成功"))
	}
}
