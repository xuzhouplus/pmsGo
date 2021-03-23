package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/lib/controller"
)

type setting struct {
	controller.App
}

var Setting = &setting{}

func (setting setting) Index(c *gin.Context) {
	result, err := model.Setting.GetPublicSettings()
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	returnData := make(map[string]interface{})
	returnData["code"] = 1
	returnData["data"] = result
	c.JSON(http.StatusOK, returnData)
}
