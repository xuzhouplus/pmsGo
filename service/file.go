package service

import (
	"context"
	"errors"
	"math"
	"pmsGo/lib/cache"
	"pmsGo/lib/config"
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/file/image"
	"pmsGo/lib/log"
	"pmsGo/lib/sync"
	"pmsGo/model"
	"pmsGo/worker"
	"strconv"
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
			if f.Poster != "" {
				f.Poster = fileLib.FullUrl(f.Poster)
			}
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

func (service File) Update(uuid string, name string, description string) (*model.File, error) {
	if uuid == "" {
		return nil, errors.New("文件id不能为空")
	}
	if name == "" {
		return nil, errors.New("文件名称不能为空")
	}
	file, err := service.FindByUuid(uuid)
	if err != nil {
		return nil, err
	}
	file.Name = name
	file.Description = description
	db := file.DB().Where("id = ?", file.ID).Updates(file)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
		return nil, db.Error
	}
	return file, nil
}

func (service File) Delete(uuid string) error {
	var one model.File
	connect := one.DB()
	connect.Where("uuid = ?", uuid)
	connect.Limit(1)
	err := connect.Find(&one).Error
	if err != nil {
		return err
	}
	result := connect.Delete(&one)
	if result.Error != nil {
		return result.Error
	}
	err = one.RemoveAll()
	if err != nil {
		return err
	}
	return nil
}

func (service File) ProcessImage(image *model.File) error {
	task, err := sync.NewTask(worker.ImageWorkerName, map[string]interface{}{
		"uuid":  image.Uuid,
		"path":  fileLib.FullPath(image.Path),
		"steps": []string{worker.ImageWorkerCreateThumbStepName, worker.ImageWorkerCreatePreviewStepName},
	})
	if err != nil {
		return err
	}
	cache.Redis.SAdd(context.TODO(), cache.Key(FileTaskCachePrefix+image.Uuid), task.UUID)
	return nil
}

func (service File) ProcessVideo(fileModel *model.File) error {
	task, err := sync.NewTask(worker.VideoWorkerName, map[string]interface{}{
		"uuid":  fileModel.Uuid,
		"path":  fileLib.FullPath(fileModel.Path),
		"steps": []string{worker.VideoWorkerCreateThumbStepName, worker.VideoWorkerCreatePosterStepName, worker.VideoWorkerCreateM3u8StepName},
	})
	if err != nil {
		return err
	}
	cache.Redis.SAdd(context.TODO(), cache.Key("file_task:"+fileModel.Uuid), task.UUID)
	cache.Redis.Expire(context.TODO(), cache.Key("file_task:"+fileModel.Uuid), time.Hour*24)
	return nil
}

func (service File) ExtractFrame(fileId interface{}, point interface{}, width interface{}, height interface{}) (map[string]string, error) {
	fileModel, err := service.FindByUuid(fileId.(string))
	if err != nil {
		return nil, err
	}
	taskJob := &worker.FrameJob{
		Path: fileLib.FullPath(fileModel.Path),
	}
	if point != nil {
		taskJob.Seek = int(point.(float64))
	}
	if width != nil {
		taskJob.Width = int(width.(float64))
	}
	if height != nil {
		taskJob.Height = int(height.(float64))
	}
	frameWorker := worker.FrameWorker{}
	result, err := frameWorker.Process(fileModel.Uuid, taskJob)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"path": result.(string),
		"url":  fileLib.FullUrl(result.(string)),
	}, nil
}

func (service File) SetPoster(fileId string, posterPath string) error {
	fileModel, err := service.FindByUuid(fileId)
	if err != nil {
		return err
	}
	fullPath := fileLib.FullPath(posterPath)
	_, err = image.Open(fullPath)
	if err != nil {
		return err
	}
	fileModel.Poster = fileLib.RelativeUrl(fileLib.PathToUrl(fileLib.Path(fullPath)))
	imageWorker := worker.ImageWorker{}
	thumbResult, err := imageWorker.Process(fileModel.Uuid, map[string]interface{}{
		"path":  fullPath,
		"steps": []string{worker.ImageWorkerCreateThumbStepName},
	})
	thumbMap := thumbResult.(map[string]string)
	fileModel.Thumb = string(fileLib.PathToUrl(fileLib.Path(thumbMap["thumb"])))
	db := fileModel.DB().Where("id = ?", fileModel.ID).Updates(fileModel)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
		return db.Error
	}
	return nil
}

func (service File) CapturePoster(fileId interface{}, point interface{}, width interface{}, height interface{}) (*model.File, error) {
	fileModel, err := service.FindByUuid(fileId.(string))
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	taskJob := &worker.FrameJob{
		Path: fileLib.FullPath(fileModel.Path),
	}
	if point != nil {
		taskJob.Seek = int(point.(float64))
	}
	if width != nil {
		taskJob.Width = int(width.(float64))
	}
	if height != nil {
		taskJob.Height = int(height.(float64))
	}
	taskJob.FileName = "poster_" + timestamp
	log.Debugf("%V", taskJob)
	frameWorker := worker.FrameWorker{}
	posterResult, err := frameWorker.Process(fileModel.Uuid, taskJob)
	if err != nil {
		return nil, err
	}
	oldPoster := fileModel.Poster
	posterPath := posterResult.(string)
	fileModel.Poster = string(fileLib.PathToUrl(fileLib.Path(posterPath)))
	imageWorker := worker.ImageWorker{}
	thumbResult, err := imageWorker.Process(fileModel.Uuid, map[string]interface{}{
		"path":  fileLib.FullPath(posterPath),
		"steps": []string{worker.ImageWorkerCreateThumbStepName},
		"uuid":  fileModel.Uuid,
		"name":  "thumb_" + timestamp,
	})
	if err != nil {
		return nil, err
	}
	oldThumb := fileModel.Thumb
	thumbMap := thumbResult.(map[string]interface{})
	fileModel.Thumb = string(fileLib.PathToUrl(fileLib.Path(thumbMap["thumb"].(string))))
	db := fileModel.DB().Where("id = ?", fileModel.ID).Updates(fileModel)
	if db.Error != nil {
		log.Errorf("%err\n", db.Error)
		return nil, db.Error
	}
	fileLib.Remove(oldPoster)
	fileLib.Remove(oldThumb)
	return fileModel, nil
}
