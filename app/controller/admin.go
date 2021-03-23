package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/lib/controller"
)

type admin struct {
	controller.App
}

var Admin = &admin{}

func (admin admin) Login(c *gin.Context) {
	requestData := make(map[string]string)
	c.BindJSON(&requestData)
	loginAdmin, error := model.Admin.Login(requestData["account"], requestData["password"])
	if error != nil {
		c.JSON(http.StatusOK, error.Error())
		return
	}
	c.JSON(http.StatusOK, loginAdmin)
}
