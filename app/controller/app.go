package controller

type app struct {
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

func (controller app) response(code int, data interface{}, message string) *response {
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
