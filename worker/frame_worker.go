package worker

const FrameWorkerName = "frame_worker"

const ExtractFrameCachePrefix = "extract_frame:"

type FrameWorker struct{}

type FrameJob struct {
	FileName string
	Path     string
	Seek     int
	Width    int
	Height   int
	Fallback func(taskId string, params interface{}, err error)
	Callback func(taskId string, params interface{}, result interface{})
}

func (FrameWorker) Process(taskId string, params interface{}) (result interface{}, err error) {
	taskParams := params.(*FrameJob)
	video, err := OpenVideo(taskParams.Path)
	if err != nil {
		return nil, err
	}
	frameParams := map[string]interface{}{
		"seek": taskParams.Seek,
	}
	if taskParams.Width == 0 {
		frameParams["width"] = video.Width
	}
	if taskParams.Height == 0 {
		frameParams["height"] = video.Height
	}
	if taskParams.FileName != "" {
		frameParams["name"] = taskParams.FileName
	}
	frame, err := createVideoFrame(taskId, video, frameParams)
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
