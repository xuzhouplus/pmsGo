package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pmsGo/app/model"
	"pmsGo/lib/controller"
)

type admin struct {
	controller.App
}

var Admin = &admin{}

func (ctl admin) Login(ctx *gin.Context) {
	requestData := make(map[string]string)
	ctx.ShouldBind(&requestData)
	loginAdmin, error := model.AdminModel.Login(requestData["account"], requestData["password"])
	if error != nil {
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeFail, nil, error.Error()))
	} else {
		session := sessions.Default(ctx)
		data, _ := json.Marshal(loginAdmin)
		session.Set("login_admin", data)
		session.Save()
		returnAttr := make(map[string]string)
		returnAttr["uuid"] = loginAdmin.Uuid
		returnAttr["type"] = loginAdmin.Type
		returnAttr["avatar"] = loginAdmin.Avatar
		returnAttr["account"] = loginAdmin.Account
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnAttr, "登录成功"))
	}
}

func (ctl admin) Auth(ctx *gin.Context) {
	returnAttr := make(map[string]interface{})
	loginAdmin := make(map[string]interface{})
	session := sessions.Default(ctx)
	loginData := session.Get("login_admin")
	if loginData == nil {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, returnAttr, "获取失败"))
		return
	}
	err := json.Unmarshal(loginData.([]byte), &loginAdmin)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, returnAttr, "获取失败"))
	}
	if loginAdmin != nil {
		returnAttr["uuid"] = loginAdmin["uuid"]
		returnAttr["type"] = loginAdmin["type"]
		returnAttr["avatar"] = loginAdmin["avatar"]
		returnAttr["account"] = loginAdmin["account"]
		ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, returnAttr, "获取成功"))
	} else {
		ctx.JSON(http.StatusBadRequest, ctl.Response(controller.CodeOk, returnAttr, "获取失败"))
	}
}

func (ctl admin) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.JSON(http.StatusOK, ctl.Response(controller.CodeOk, nil, "退出成功"))
}
