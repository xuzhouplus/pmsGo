package worker

import (
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/file/image"
	"pmsGo/lib/log"
	"pmsGo/model"
	"strconv"
)

const CarouselWorkerName = "carousel_worker"

type CarouselWorker struct{}

func (CarouselWorker) Process(taskId string, params interface{}) (result interface{}, err error) {
	taskParams := params.(map[string]interface{})
	image, err := openCarousel(taskParams["path"].(string))
	if err != nil {
		return nil, err
	}
	carousel, err := createCarousel(image, taskParams)
	if err != nil {
		return nil, err
	}
	return carousel, nil
}

func (CarouselWorker) Fallback(taskId string, params interface{}, err error) {
	taskParams := params.(map[string]interface{})
	videoModel := &model.Carousel{
		Status: model.FileStatusError,
		Error:  err.Error(),
	}
	db := videoModel.DB().Where("uuid = ?", taskParams["uuid"]).Updates(videoModel)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
	}
}

func (CarouselWorker) Callback(taskId string, params interface{}, result interface{}) {
	taskParams := params.(map[string]interface{})
	taskResult := result.(map[string]interface{})
	carouselModel := &model.Carousel{
		Width:  taskResult["width"].(int),
		Height: taskResult["height"].(int),
		Url:    taskResult["url"].(string),
		Thumb:  taskResult["thumb"].(string),
		Status: model.CarouselStatusEnabled,
		Error:  "",
	}
	db := carouselModel.DB().Where("uuid = ?", taskParams["uuid"]).Updates(carouselModel)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
	}
}

func openCarousel(path string) (*image.Image, error) {
	return image.Open(path)
}

func createCarousel(image *image.Image, params map[string]interface{}) (map[string]interface{}, error) {
	carouselWidth := 1920
	if params["carousel_width"] != nil {
		carouselWidth = params["carousel_width"].(int)
	}
	carouselHeight := 1080
	if params["carousel_height"] != nil {
		carouselHeight = params["carousel_height"].(int)
	}
	carouselExtension := "jpg"
	if params["carousel_extension"] != nil {
		carouselExtension = params["carousel_extension"].(string)
	}
	carouselName := params["uuid"].(string) + "_" + strconv.Itoa(carouselWidth) + "_" + strconv.Itoa(carouselHeight)
	carouselFile, err := image.CreateCarousel(carouselName, carouselWidth, carouselHeight, carouselExtension)
	if err != nil {
		log.Errorf("%err\n", err)
		return nil, err
	}
	thumbWidth := 320
	if params["thumb_width"] != nil {
		thumbWidth = params["thumb_width"].(int)
	}
	thumbHeight := 180
	if params["thumb_height"] != nil {
		thumbHeight = params["thumb_height"].(int)
	}
	thumbExtension := "jpg"
	if params["thumb_extension"] != nil {
		thumbExtension = params["thumb_extension"].(string)
	}
	thumbName := params["uuid"].(string) + "_" + strconv.Itoa(thumbWidth) + "_" + strconv.Itoa(thumbHeight)
	thumb, err := carouselFile.CreateResize(thumbName, thumbWidth, thumbHeight, thumbExtension)
	if err != nil {
		log.Errorf("%err\n", err)
		return nil, err
	}
	return map[string]interface{}{
		"width":  carouselWidth,
		"height": carouselHeight,
		"url":    fileLib.RelativePath(fileLib.Path(carouselFile.FullPath())),
		"thumb":  fileLib.RelativePath(fileLib.Path(thumb.FullPath())),
	}, nil
}
