package service

import (
	"context"
	"errors"
	"math"
	"pmsGo/lib/cache"
	"pmsGo/lib/config"
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/sync"
	"pmsGo/model"
	"pmsGo/worker"
	"time"
)

const (
	FileTaskCachePrefix    = "file_task:"
	TaskProcessCachePrefix = "task_process:"
)

type File struct {
}

var FileService = &File{}

func (service File) GetMimeTypes(fileType string) []string {
	if fileType == "" {
		return nil
	}
	fileTypes := make([]string, 0)
	switch fileType {
	case fileLib.TypeVideo:
		fileTypes = config.Config.Web.Upload.Video.Extensions
	case fileLib.TypeImage:
		fileTypes = config.Config.Web.Upload.Image.Extensions
	default:
		return append(fileTypes, fileType)
	}
	if len(fileTypes) == 1 && fileTypes[0] == "*" {
		return nil
	}
	return fileTypes
}

func (service File) List(page int, limit int, fields []string, fileType string, name string) (map[string]interface{}, error) {
	var files []model.File
	fileModel := &model.File{}
	connect := fileModel.DB()
	if len(fields) > 0 {
		connect.Select(fields)
	}
	if name != "" {
		connect.Where("name like ?", "%"+name+"%")
	}
	if fileType != "" {
		connect.Where("type", fileType)
	}
	if page < 0 {
		page = 0
	}
	if limit == 0 {
		limit = 10
	}
	connect.Offset(page * limit)
	connect.Limit(limit)
	returnData := make(map[string]interface{})
	if connect.Find(&files).Error != nil {
		return returnData, errors.New("获取文件列表失败")
	}
	if len(files) > 0 {
		for i, f := range files {
			f.Path = fileLib.FullUrl(f.Path)
			f.Thumb = fileLib.FullUrl(f.Thumb)
			f.Preview = fileLib.FullUrl(f.Preview)
			files[i] = f
		}
	}
	returnData["files"] = files
	returnData["size"] = limit
	returnData["page"] = page
	var total int64
	connect.Offset(-1)
	connect.Limit(-1)
	connect.Count(&total)
	returnData["total"] = total
	returnData["count"] = math.Ceil(float64(total) / float64(limit))
	return returnData, nil
}

func (service File) Upload(uploaded *fileLib.Upload, name string, description string) (*model.File, error) {
	fileModel := &model.File{}
	fileModel.Uuid = uploaded.Uuid
	fileModel.Name = name
	fileModel.Description = description
	fileModel.Path = fileLib.RelativePath(uploaded.Path())
	fileModel.Type = uploaded.FileType
	fileModel.Extension = uploaded.Extension
	fileModel.Status = model.FileStatusUploaded
	connect := fileModel.DB()
	result := connect.Create(&fileModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return fileModel, nil
}

func (service File) FindOne(id int) (*model.File, error) {
	one := &model.File{}
	connect := one.DB()
	connect.Where("id = ?", id)
	connect.Limit(1)
	err := connect.Find(&one).Error
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (service File) FindByUuid(uuid string) (*model.File, error) {
	one := &model.File{}
	connect := one.DB()
	connect.Where("uuid = ?", uuid)
	connect.Limit(1)
	err := connect.Find(&one).Error
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (service File) Delete(id int) error {
	var one model.File
	connect := one.DB()
	connect.Where("id = ?", id)
	connect.Limit(1)
	err := connect.Find(&one).Error
	if err != nil {
		return err
	}
	result := connect.Delete(&one)
	if result.Error != nil {
		return result.Error
	}
	err = one.RemoveFile()
	if err != nil {
		return err
	}
	err = one.RemoveThumb()
	if err != nil {
		return err
	}
	err = one.RemovePreview()
	if err != nil {
		return err
	}
	err = one.RemoveDir()
	if err != nil {
		return err
	}
	return nil
}

func (service File) ProcessImage(image *model.File) error {
	task, err := sync.NewTask(worker.ImageWorkerName, map[string]interface{}{
		"id":    image.ID,
		"uuid":  image.Uuid,
		"path":  fileLib.FullPath(image.Path),
		"steps": []string{worker.ImageWorkerCreateThumbStepName, worker.ImageWorkerCreatePreviewStepName},
	})
	if err != nil {
		return err
	}
	cache.Redis.SAdd(context.TODO(), FileTaskCachePrefix+image.Uuid, task.UUID)
	return nil
}

func (service File) ProcessVideo(fileModel *model.File) error {
	task, err := sync.NewTask(worker.VideoWorkerName, map[string]interface{}{
		"id":    fileModel.ID,
		"path":  fileLib.FullPath(fileModel.Path),
		"steps": []string{worker.VideoWorkerCreateThumbStepName, worker.VideoWorkerCreatePosterStepName, worker.VideoWorkerCreateM3u8StepName},
	})
	if err != nil {
		return err
	}
	cache.Redis.SAdd(context.TODO(), "file_task:"+fileModel.Uuid, task.UUID)
	cache.Redis.Expire(context.TODO(), "file_task:"+fileModel.Uuid, time.Hour*24)
	return nil
}

func (service File) ExtractFrame(fileId string, point int64, width int64, height int64) (string, error) {
	fileModel, err := service.FindByUuid(fileId)
	if err != nil {
		return "", err
	}
	taskJob := &worker.FrameJob{
		Path:     fileLib.FullPath(fileModel.Path),
		Seek:     int(point),
		Width:    int(width),
		Height:   int(height),
		Callback: ExtractFrameCallback,
		Fallback: ExtractFrameFallback,
	}
	task, err := sync.NewTask(worker.FrameWorkerName, taskJob)
	if err != nil {
		return "", err
	}
	return task.UUID, nil
}

const ExtractFrameCachePrefix = "extract_frame:"

func ExtractFrameFallback(taskId string, params interface{}, err error) {
	cache.Redis.Set(context.TODO(), ExtractFrameCachePrefix+taskId, map[string]string{
		"status": "fail",
		"error":  err.Error(),
	}, time.Hour*24)
}

func ExtractFrameCallback(taskId string, params interface{}, result interface{}) {
	cache.Redis.Set(context.TODO(), ExtractFrameCachePrefix+taskId, map[string]interface{}{
		"status": "fail",
		"result": result,
	}, time.Hour*24)
}

func (service File) GetFrame(taskId string) interface{} {
	return cache.Redis.Get(context.TODO(), ExtractFrameCachePrefix+taskId)
}
