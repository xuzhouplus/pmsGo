package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
)

type admin struct {
	app
}

var Admin = &admin{}

func (admin admin) Login(c *gin.Context) {
	requestData := make(map[string]string)
	c.BindJSON(&requestData)
	loginAdmin, error := model.Admin.Login(requestData["account"], requestData["password"])
	if error != nil {
		c.JSON(http.StatusOK, admin.response(CodeOk, nil, error.Error()))
	} else {
		c.JSON(http.StatusOK, admin.response(CodeOk, loginAdmin, "登录成功"))
	}
}
