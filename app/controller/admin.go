package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/app/model"
)

type admin struct {
}

var Admin = &admin{}

func (admin admin) Login(c *gin.Context) {
	fmt.Println(c.Request.PostFormValue("account"))
	loginForm := &model.Login{}
	c.BindJSON(loginForm)
	c.JSON(http.StatusOK, loginForm)
}
