package worker

import (
	"fmt"
	"github.com/floostack/transcoder"
	"math"
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/file/video"
	"pmsGo/lib/helper"
	"pmsGo/lib/log"
	"pmsGo/model"
	"strconv"
)

const VideoWorkerName = "video_worker"
const (
	VideoWorkerCreateThumbStepName  = "video_worker_create_thumb"
	VideoWorkerCreatePosterStepName = "video_worker_create_poster"
	VideoWorkerCreateM3u8StepName   = "video_worker_create_m3u8"
)

type VideoWorker struct {
}

func (VideoWorker) Process(taskId string, params interface{}) (result interface{}, err error) {
	taskParams := params.(map[string]interface{})
	video, err := OpenVideo(taskParams["path"].(string))
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	data["width"] = video.Width
	data["height"] = video.Height
	taskStep := TaskSteps(taskParams["steps"])
	_, isIn := helper.IsInSlice(taskStep, VideoWorkerCreateThumbStepName)
	if isIn {
		thumb, err := createVideoThumb(taskId, video, taskParams)
		if err != nil {
			return nil, err
		}
		data["thumb"] = thumb
	}
	_, isIn = helper.IsInSlice(taskStep, VideoWorkerCreatePosterStepName)
	if isIn {
		poster, err := createVideoPoster(taskId, video, taskParams)
		if err != nil {
			return nil, err
		}
		data["poster"] = poster
	}
	_, isIn = helper.IsInSlice(taskStep, VideoWorkerCreateM3u8StepName)
	if isIn {
		m3u8, err := createVideoM3u8(taskId, video, taskParams)
		if err != nil {
			return nil, err
		}
		data["m3u8"] = m3u8
	}
	return data, nil
}

func (VideoWorker) Fallback(taskId string, params interface{}, err error) {
	taskParams := params.(map[string]interface{})
	videoModel := &model.File{
		Status: model.FileStatusError,
		Error:  err.Error(),
	}
	db := videoModel.DB().Where("uuid = ?", taskParams["uuid"]).Updates(videoModel)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
	}
	ClearTaskProcessStatus(taskId)
}

func (VideoWorker) Callback(taskId string, params interface{}, result interface{}) {
	taskParams := params.(map[string]interface{})
	taskResult := result.(map[string]interface{})
	videoModel := &model.File{
		Width:  taskResult["width"].(int),
		Height: taskResult["height"].(int),
		Status: model.FileStatusEnabled,
		Error:  "",
	}
	if taskResult["thumb"] != nil {
		videoModel.Thumb = taskResult["thumb"].(string)
	}
	if taskResult["poster"] != nil {
		videoModel.Poster = taskResult["poster"].(string)
	}
	if taskResult["m3u8"] != nil {
		videoModel.Preview = taskResult["m3u8"].(string)
	}
	db := videoModel.DB().Where("uuid = ?", taskParams["uuid"]).Updates(videoModel)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
	}
	ClearTaskProcessStatus(taskId)
}

func OpenVideo(path string) (*video.Video, error) {
	return video.Open(fileLib.FullPath(path))
}

func extractVideoFrame(video *video.Video, param map[string]interface{}) (string, <-chan transcoder.Progress, error) {
	duration, err := strconv.ParseFloat(video.Duration, 64)
	if err != nil {
		log.Errorf("%err\n", err)
		return "", nil, err
	} else {
		seek := 0.0
		switch param["seek"].(type) {
		case string:
			seek, err = strconv.ParseFloat(param["seek"].(string), 64)
			if err != nil {
				log.Errorf("%err\n", err)
				return "", nil, err
			}
			seek = math.Max(seek, 0.0)
		case int:
			seek = math.Max(float64(param["seek"].(int)), seek)
		case float64:
			seek = math.Max(param["seek"].(float64), 0.0)
		case int64:
			seek = math.Max(float64(param["seek"].(int64)), 0.0)
		default:
			return "", nil, fmt.Errorf("seek类型错误:%t", param["seek"])
		}
		seek = math.Min(seek, duration)
		width := 320
		height := 180
		if param["width"] != nil {
			width = param["width"].(int)
		}
		if param["height"] != nil {
			height = param["height"].(int)
		}
		name := ""
		if param["name"] != nil {
			name = param["name"].(string)
		}
		path, progressChannel, err := video.CreateThumb(name, width, height, "jpg", helper.SecondToTime(seek))
		if err != nil {
			log.Errorf("%err\n", err)
			return path, progressChannel, err
		} else {
			return path, progressChannel, nil
		}
	}
}

func createVideoFrame(taskId string, video *video.Video, param map[string]interface{}) (string, error) {
	log.Debugf("video frame sync task:%v\n", param)
	path, progressChannel, err := extractVideoFrame(video, map[string]interface{}{
		"name":   param["name"].(string),
		"seek":   param["seek"].(int),
		"height": param["height"].(int),
		"width":  param["width"].(int),
	})
	if err != nil {
		log.Errorf("%err\n", err)
		status := map[string]interface{}{"error": err.Error(), "status": "fail"}
		SetTaskProcessStatus(taskId, "createVideoFrame", status)
		return "", err
	}
	for progress := range progressChannel {
		current := int(progress.GetProgress() * 100)
		status := map[string]interface{}{"error": "", "status": "progress", "progress": current}
		SetTaskProcessStatus(taskId, "createVideoFrame", status)
	}
	return fileLib.RelativePath(fileLib.Path(path)), nil
}

func createVideoThumb(taskId string, video *video.Video, param map[string]interface{}) (string, error) {
	log.Debugf("video thumb sync task:%v\n", param)
	path, progressChannel, err := extractVideoFrame(video, map[string]interface{}{
		"seek":   3.0,
		"height": 180,
		"width":  320,
	})
	if err != nil {
		log.Errorf("%err\n", err)
		status := map[string]interface{}{"error": err.Error(), "status": "fail"}
		SetTaskProcessStatus(taskId, "createVideoThumb", status)
		return "", err
	}
	for progress := range progressChannel {
		current := int(progress.GetProgress() * 100)
		status := map[string]interface{}{"error": "", "status": "progress", "progress": current}
		SetTaskProcessStatus(taskId, "createVideoThumb", status)
	}
	return fileLib.RelativePath(fileLib.Path(path)), nil
}

func createVideoPoster(taskId string, video *video.Video, param map[string]interface{}) (string, error) {
	log.Debugf("video poster sync task:%v\n", param)
	path, progressChannel, err := extractVideoFrame(video, map[string]interface{}{
		"seek":   3.0,
		"height": video.Height,
		"width":  video.Width,
	})
	if err != nil {
		log.Errorf("%err\n", err)
		status := map[string]interface{}{"error": err.Error(), "status": "fail"}
		SetTaskProcessStatus(taskId, "createVideoPoster", status)
		return "", err
	}
	for progress := range progressChannel {
		current := int(progress.GetProgress() * 100)
		status := map[string]interface{}{"error": "", "status": "progress", "progress": current}
		SetTaskProcessStatus(taskId, "createVideoPoster", status)
	}
	return fileLib.RelativePath(fileLib.Path(path)), nil
}

func createVideoM3u8(taskId string, video *video.Video, param map[string]interface{}) (string, error) {
	log.Debugf("video m3u8 sync task:%v\n", param)
	width := 720
	if param["width"] != nil {
		width = param["width"].(int)
	}
	height := 576
	if param["height"] != nil {
		height = param["height"].(int)
	}
	path, progressChannel, err := video.CreateM3u8(width, height)
	if err != nil {
		log.Errorf("%err\n", err)
		status := map[string]interface{}{"error": err.Error(), "status": "fail"}
		SetTaskProcessStatus(taskId, "createVideoM3u8", status)
		return "", err
	}
	for progress := range progressChannel {
		current := int(progress.GetProgress() * 100)
		status := map[string]interface{}{"error": "", "status": "progress", "progress": current}
		SetTaskProcessStatus(taskId, "createVideoM3u8", status)
	}
	status := map[string]interface{}{"error": "", "status": "success", "progress": 100}
	SetTaskProcessStatus(taskId, "createVideoM3u8", status)
	return fileLib.RelativePath(fileLib.Path(path)), nil
}
