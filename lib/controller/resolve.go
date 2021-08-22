package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pmsGo/lib/helper"
	"pmsGo/lib/log"
	"reflect"
)

type action struct {
	Name          string
	Verbs         []string
	Authenticator string
	Handler       interface{}
}

func (a action) Func() gin.HandlerFunc {
	switch a.Handler.(type) {
	case gin.HandlerFunc:
		return a.Handler.(gin.HandlerFunc)
	case reflect.Method:
		handlerFunc := a.Handler.(reflect.Method)
		return func(context *gin.Context) {
			result := handlerFunc.Func.Call([]reflect.Value{reflect.ValueOf(context)})
			log.Debug(result)
		}
	default:
		return func(context *gin.Context) {
			context.JSON(http.StatusNotFound, nil)
		}
	}
}

type Resolve struct {
	Controller    reflect.Type
	Verbs         map[string][]string
	Authenticator Authenticator
	Actions       map[string]*action
}

func NewResolve(controller AppInterface) *Resolve {
	resolve := &Resolve{}
	resolve.Actions = make(map[string]*action)
	resolve.Verbs = make(map[string][]string)
	resolve.ReflectController(controller)
	resolve.Verbs = controller.Verbs()
	resolve.Authenticator = controller.Authenticator()
	resolve.ReflectActions(controller)
	return resolve
}
func (r *Resolve) ReflectController(controller AppInterface) {
	r.Controller = reflect.TypeOf(controller).Elem()
}
func (r *Resolve) ReflectActions(controller AppInterface) {
	methodNum := r.Controller.NumMethod()
	actions := controller.Actions()
	if actions != nil {
		for name, handlerFunc := range actions {
			r.Actions[name] = &action{
				Name:          name,
				Handler:       handlerFunc,
				Verbs:         r.ReflectVerb(name),
				Authenticator: r.ReflectAuthenticator(name),
			}
		}
	}
	for loopIndex := 0; loopIndex < methodNum; loopIndex++ {
		method := r.Controller.Method(loopIndex)
		if method.IsExported() && method.Type.String() == "func(controller.admin, *gin.Context)" {
			actionReflect := &action{
				Name:          method.Name,
				Handler:       method,
				Verbs:         r.ReflectVerb(method.Name),
				Authenticator: r.ReflectAuthenticator(method.Name),
			}
			r.Actions[method.Name] = actionReflect
		}
	}
}

func (r Resolve) ReflectVerb(action string) []string {
	if r.Verbs == nil {
		return []string{
			Any,
		}
	}
	verbs := r.Verbs[action]
	if verbs == nil {
		return []string{
			Any,
		}
	}
	_, result := helper.IsInSlice(verbs, Any)
	if result {
		return []string{
			Any,
		}
	}
	return verbs
}

func (r Resolve) ReflectAuthenticator(action string) string {
	if r.Authenticator.Excepts != nil {
		_, result := helper.IsInSlice(r.Authenticator.Excepts, action)
		if result {
			return Except
		}
	}
	if r.Authenticator.Optionals == nil {
		_, result := helper.IsInSlice(r.Authenticator.Optionals, action)
		if result {
			return Optional
		}
	}
	return Forbidden
}
