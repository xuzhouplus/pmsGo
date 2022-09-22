package worker

import (
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/file/image"
	"pmsGo/lib/helper"
	"pmsGo/lib/log"
	"pmsGo/model"
	"strconv"
)

const ImageWorkerName = "image_worker"
const (
	ImageWorkerCreateThumbStepName   = "image_worker_create_thumb"
	ImageWorkerCreatePreviewStepName = "image_worker_create_preview"
)

type ImageWorker struct{}

func (ImageWorker) Process(taskId string, params interface{}) (result interface{}, err error) {
	taskParams := params.(map[string]interface{})
	image, err := openImage(taskParams["path"].(string))
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	data["width"] = image.Width
	data["height"] = image.Height
	taskStep := TaskSteps(taskParams["steps"])
	_, isIn := helper.IsInSlice(taskStep, ImageWorkerCreateThumbStepName)
	if isIn {
		thumb, err := createImageResize(image, taskParams)
		if err != nil {
			return nil, err
		}
		data["thumb"] = thumb
	}
	_, isIn = helper.IsInSlice(taskStep, ImageWorkerCreatePreviewStepName)
	if isIn {
		preview, err := createImageCompress(image, taskParams)
		if err != nil {
			return nil, err
		}
		data["preview"] = preview
	}
	return data, nil
}

func (ImageWorker) Fallback(taskId string, params interface{}, err error) {
	taskParams := params.(map[string]interface{})
	videoModel := &model.File{
		Status: model.FileStatusError,
		Error:  err.Error(),
	}
	db := videoModel.DB().Where("uuid = ?", taskParams["uuid"]).Updates(videoModel)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
	}
}

func (ImageWorker) Callback(taskId string, params interface{}, result interface{}) {
	taskParams := params.(map[string]interface{})
	taskResult := result.(map[string]interface{})
	fileModel := &model.File{
		Width:  taskResult["width"].(int),
		Height: taskResult["height"].(int),
		Status: model.FileStatusEnabled,
		Error:  "",
	}
	if taskResult["thumb"] != nil {
		fileModel.Thumb = taskResult["thumb"].(string)
	}
	if taskResult["preview"] != nil {
		fileModel.Preview = taskResult["preview"].(string)
	}
	db := fileModel.DB().Where("uuid = ?", taskParams["uuid"]).Updates(fileModel)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
	}
}

func openImage(path string) (*image.Image, error) {
	return image.Open(path)
}

func createImageResize(image *image.Image, params map[string]interface{}) (string, error) {
	width := 320
	if params["width"] != nil {
		width = params["width"].(int)
	}
	height := 180
	if params["height"] != nil {
		height = params["height"].(int)
	}
	extension := ""
	if params["extension"] != nil {
		extension = params["extension"].(string)
	}
	name := ""
	if params["name"] == nil {
		name = params["uuid"].(string) + "_" + strconv.Itoa(width) + "_" + strconv.Itoa(height)
	} else {
		name = params["name"].(string)
	}
	thumb, err := image.CreateResize(name, width, height, extension)
	if err != nil {
		log.Errorf("%err\n", err)
		return "", err
	}
	return fileLib.RelativePath(fileLib.Path(thumb.FullPath())), nil
}

func createImageCompress(image *image.Image, params map[string]interface{}) (string, error) {
	quality := 62
	if params["quality"] != nil {
		quality = params["quality"].(int)
	}
	name := ""
	if params["name"] == nil {
		name = params["uuid"].(string) + "_" + strconv.Itoa(image.Width) + "_" + strconv.Itoa(image.Height) + "_" + strconv.Itoa(quality)
	} else {
		name = params["name"].(string)
	}
	thumb, err := image.CreateCompress(name, quality)
	if err != nil {
		log.Errorf("%err\n", err)
		return "", err
	}
	return fileLib.RelativePath(fileLib.Path(thumb.FullPath())), nil
}
