package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Method string

const (
	Post    = http.MethodPost
	Put     = http.MethodPut
	Options = http.MethodOptions
	Get     = http.MethodGet
	Delete  = http.MethodDelete
	Head    = http.MethodHead
	Patch   = http.MethodPatch
	Trace   = http.MethodTrace
	Connect = http.MethodConnect
	Any     = "*"
)

// Except 不需要登录
const Except = "except"

// Optional 可以不登录
const Optional = "optional"

// Forbidden 拒绝访问
const Forbidden = "forbidden"

type Authenticator struct {
	Excepts   []string
	Optionals []string
}

type AppInterface interface {
	Verbs() map[string][]string          // Verbs 配置方法请求方式
	Authenticator() Authenticator        // Authenticator 配置方法登录限制
	Actions() map[string]gin.HandlerFunc // Actions 配置方法映射
}

type AppController struct {
	Excepts   []string
	Optionals []string
}

const (
	CodeOk   = 1
	CodeFail = 0
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (controller AppController) Verbs() map[string][]string {
	return nil
}

func (controller AppController) Authenticator() Authenticator {
	authenticator := Authenticator{
		controller.Excepts,
		controller.Optionals,
	}
	return authenticator
}

func (controller AppController) Actions() map[string]gin.HandlerFunc {
	return nil
}

func (controller AppController) CodeOk() int {
	return CodeOk
}

func (controller AppController) CodeFail() int {
	return CodeFail
}

func (controller AppController) Response(code int, data interface{}, message string) *response {
	returnData := &response{}
	if code == 0 {
		returnData.Code = CodeFail
	} else {
		returnData.Code = CodeOk
	}
	returnData.Data = data
	returnData.Message = message
	return returnData
}

func (controller AppController) ResponseOk(data interface{}, message string) *response {
	returnData := &response{}
	returnData.Code = CodeOk
	returnData.Data = data
	returnData.Message = message
	return returnData
}

func (controller AppController) ResponseFail(data interface{}, message string) *response {
	returnData := &response{}
	returnData.Code = CodeFail
	returnData.Data = data
	returnData.Message = message
	return returnData
}
