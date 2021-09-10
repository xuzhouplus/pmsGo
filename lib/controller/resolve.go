package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/lib/helper"
	"reflect"
)

//控制器请求方法
type action struct {
	Name          string      //方法名称
	Verbs         []string    //方法接受的请求方式
	Authenticator string      //方法登录授权验证方式
	Handler       interface{} //请求处理方法
}

// Resolve 控制器解析结果
type Resolve struct {
	Value   reflect.Value      //控制器实例
	Type    reflect.Type       //控制器类型
	Actions map[string]*action //控制器中的方法
}

// NewResolve 获取控制器解析结果
func NewResolve(controller AppInterface) *Resolve {
	resolve := &Resolve{
		Actions: make(map[string]*action),
	}
	resolve.ReflectController(controller)
	resolve.ReflectActions(controller)
	return resolve
}

// Handle 获取请求处理方法
func (r Resolve) Handle(method string) gin.HandlerFunc {
	action := r.Actions[method]
	switch action.Handler.(type) {
	case gin.HandlerFunc:
		return action.Handler.(gin.HandlerFunc)
	case reflect.Method:
		handlerFunc := action.Handler.(reflect.Method)
		return func(context *gin.Context) {
			in := make([]reflect.Value, 2)
			in[0] = r.Value
			in[1] = reflect.ValueOf(context)
			handlerFunc.Func.Call(in)
		}
	default:
		return func(context *gin.Context) {
			context.JSON(http.StatusNotImplemented, nil)
		}
	}
}

// GetControllerName 获取控制器名称
func (r Resolve) GetControllerName() string {
	return helper.CamelToLine(r.Type.Elem().Name())
}

// GetActions 获取控制器方法
func (r Resolve) GetActions() map[string]*action {
	return r.Actions
}

// GetAction 获取指定的控制器方法
func (r Resolve) GetAction(method string) *action {
	return r.Actions[method]
}

// ReflectController 反射控制器类型和实例
func (r *Resolve) ReflectController(controller AppInterface) {
	r.Type = reflect.TypeOf(controller)
	r.Value = reflect.ValueOf(controller)
}

// ReflectVerbs 解析方法请求方式配置映射到方法名称上
func (r Resolve) ReflectVerbs(controller AppInterface) map[string][]string {
	methodVerbs := make(map[string][]string)
	actionVerbs := controller.Verbs()
	for action, verbs := range actionVerbs {
		method := r.ResolveAction(action)
		methodVerbs[method] = verbs
	}
	return methodVerbs
}

// ReflectAuthenticator 解析方法登录授权配置
func (r Resolve) ReflectAuthenticator(controller AppInterface) Authenticator {
	methodAuthenticator := Authenticator{
		Excepts:   []string{},
		Optionals: []string{},
	}
	authenticator := controller.Authenticator()
	if authenticator.Excepts != nil {
		for _, exceptAction := range authenticator.Excepts {
			methodAuthenticator.Excepts = append(methodAuthenticator.Excepts, r.ResolveAction(exceptAction))
		}
	}
	if authenticator.Optionals != nil {
		for _, optionalAction := range authenticator.Optionals {
			methodAuthenticator.Optionals = append(methodAuthenticator.Optionals, r.ResolveAction(optionalAction))
		}
	}
	return methodAuthenticator
}

// ReflectActions 反射控制器方法
func (r *Resolve) ReflectActions(controller AppInterface) {
	//获取控制器方法数量
	methodNum := r.Type.NumMethod()
	//获取控制器方法请求方式配置
	verbs := r.ReflectVerbs(controller)
	//获取控制器方法登录权限控制
	authenticator := r.ReflectAuthenticator(controller)
	//获取控制器Actions方法中定义的方法
	actions := controller.Actions()
	if actions != nil {
		for name, handlerFunc := range actions {
			name = helper.CamelToLine(name)
			r.Actions[name] = &action{
				Name:          name,
				Handler:       handlerFunc,
				Verbs:         r.ResolveVerb(verbs, name),
				Authenticator: r.ResolveAuthenticator(authenticator, name),
			}
		}
	}
	//反射获取控制器中的方法
	methodType := "func(*controller." + r.Type.Elem().Name() + ", *gin.Context)"
	for loopIndex := 0; loopIndex < methodNum; loopIndex++ {
		//反射获取方法信息
		method := r.Type.Method(loopIndex)
		//判断方法是否公开和方法类型
		if method.IsExported() && method.Type.String() == methodType {
			//把方法名转换为中横线，方便授权匹配
			name := helper.CamelToLine(method.Name)
			r.Actions[name] = &action{
				Name:          name,
				Handler:       method,
				Verbs:         r.ResolveVerb(verbs, name),
				Authenticator: r.ResolveAuthenticator(authenticator, name),
			}
		}
	}
}

// ResolveAction 获取方法中横线名称
func (r Resolve) ResolveAction(action string) string {
	return helper.CamelToLine(action)
}

// ResolveVerb 获取方法请求方式
func (r Resolve) ResolveVerb(verbs map[string][]string, action string) []string {
	if verbs == nil {
		return []string{
			Any,
		}
	}
	actionVerbs := verbs[action]
	if actionVerbs == nil {
		return []string{
			Any,
		}
	}
	_, result := helper.IsInSlice(actionVerbs, Any)
	if result {
		return []string{
			Any,
		}
	}
	return actionVerbs
}

// ResolveAuthenticator 获取方法登录权限配置
func (r Resolve) ResolveAuthenticator(authenticator Authenticator, action string) string {
	if len(authenticator.Excepts) > 0 {
		_, result := helper.IsInSlice(authenticator.Excepts, action)
		if result {
			return Except
		}
	}
	if len(authenticator.Optionals) > 0 {
		_, result := helper.IsInSlice(authenticator.Optionals, action)
		if result {
			return Optional
		}
	}
	return Forbidden
}
