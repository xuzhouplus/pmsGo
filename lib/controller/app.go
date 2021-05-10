package controller

type App struct {
	Except   []string
	Optional []string
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

func (controller App) CodeOk() int {
	return CodeOk
}
func (controller App) CodeFail() int {
	return CodeFail
}
func (controller App) Response(code int, data interface{}, message string) *response {
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

func (controller App) ResponseOk(data interface{}, message string) *response {
	returnData := &response{}
	returnData.Code = CodeOk
	returnData.Data = data
	returnData.Message = message
	return returnData
}

func (controller App) ResponseFail(data interface{}, message string) *response {
	returnData := &response{}
	returnData.Code = CodeFail
	returnData.Data = data
	returnData.Message = message
	return returnData
}
