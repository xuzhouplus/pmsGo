package worker

const FrameWorkerName = "frame_worker"

const ExtractFrameCachePrefix = "extract_frame:"

type FrameWorker struct{}

type FrameJob struct {
	Path     string
	Seek     int
	Width    int
	Height   int
	Fallback func(taskId string, params interface{}, err error)
	Callback func(taskId string, params interface{}, result interface{})
}

func (FrameWorker) Process(taskId string, params interface{}) (result interface{}, err error) {
	taskParams := params.(FrameJob)
	video, err := OpenVideo(taskParams.Path)
	if err != nil {
		return nil, err
	}
	frame, err := createVideoFrame(taskId, video, map[string]interface{}{
		"seek":   taskParams.Seek,
		"width":  taskParams.Width,
		"height": taskParams.Height,
	})
	if err != nil {
		return nil, err
	}
	return frame, nil
}

func (FrameWorker) Fallback(taskId string, params interface{}, err error) {
	taskParams := params.(FrameJob)
	fallback := taskParams.Fallback
	if fallback != nil {
		fallback(taskId, params, err)
	}
}

func (FrameWorker) Callback(taskId string, params interface{}, result interface{}) {
	taskParams := params.(FrameJob)
	callback := taskParams.Callback
	if callback != nil {
		callback(taskId, params, result)
	}
}
