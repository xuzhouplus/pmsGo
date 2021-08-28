package controller

import (
	"fmt"
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

type Resolve struct {
	Value   reflect.Value
	Type    reflect.Type
	Actions map[string]*action
}

func NewResolve(controller AppInterface) *Resolve {
	resolve := &Resolve{
		Actions: make(map[string]*action),
	}
	resolve.ReflectController(controller)
	resolve.ReflectActions(controller)
	return resolve
}
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
			result := handlerFunc.Func.Call(in)
			log.Debug(result)
		}
	default:
		fmt.Println(action)
		return func(context *gin.Context) {
			context.JSON(http.StatusNotFound, nil)
		}
	}
}
func (r Resolve) GetControllerName() string {
	return r.Type.Elem().Name()
}
func (r Resolve) GetActions() map[string]*action {
	return r.Actions
}
func (r Resolve) GetAction(method string) *action {
	return r.Actions[method]
}
func (r *Resolve) ReflectController(controller AppInterface) {
	r.Type = reflect.TypeOf(controller)
	r.Value = reflect.ValueOf(controller)
}
func (r Resolve) ReflectVerbs(controller AppInterface) map[string][]string {
	methodVerbs := make(map[string][]string)
	actionVerbs := controller.Verbs()
	for action, verbs := range actionVerbs {
		method := r.ResolveAction(action)
		methodVerbs[method] = verbs
	}
	return methodVerbs
}
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
func (r *Resolve) ReflectActions(controller AppInterface) {
	methodNum := r.Type.NumMethod()
	verbs := r.ReflectVerbs(controller)
	authenticator := r.ReflectAuthenticator(controller)
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
	for loopIndex := 0; loopIndex < methodNum; loopIndex++ {
		method := r.Type.Method(loopIndex)
		if method.IsExported() && method.Type.String() == ("func(*controller."+r.GetControllerName()+", *gin.Context)") {
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
func (r Resolve) ResolveAction(action string) string {
	return helper.CamelToLine(action)
}
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

func (r Resolve) ResolveAuthenticator(authenticator Authenticator, action string) string {
	if len(authenticator.Excepts) > 0 {
		_, result := helper.IsInSlice(authenticator.Excepts, action)
		if result {
			return Except
		}
	}
	if len(authenticator.Optionals) > 0 {
		fmt.Println(authenticator.Optionals)
		fmt.Println(action)
		_, result := helper.IsInSlice(authenticator.Optionals, action)
		if result {
			return Optional
		}
	}
	return Forbidden
}
