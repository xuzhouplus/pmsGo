package sync

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
)

type Worker interface {
	Process(taskId string, params interface{}) (result interface{}, err error)
	Fallback(taskId string, params interface{}, err error)
	Callback(taskId string, params interface{}, result interface{})
}
