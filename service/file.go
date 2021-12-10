package service

import (
	"errors"
	"math"
	"pmsGo/lib/cache"
	"pmsGo/lib/config"
	fileLib "pmsGo/lib/file"
	imageLib "pmsGo/lib/file/image"
	"pmsGo/lib/file/video"
	"pmsGo/lib/log"
	"pmsGo/lib/sync"
	"pmsGo/model"
	"strconv"
)

const (
	VideoProgressCacheKey       = "progress:video"
	ImageProgressCacheKey       = "progress:image"
	ProgressCacheTtl            = 60 * 60 * 24
	CreateImageThumbSyncTaskKey = "CreateImageThumb"
	GetVideoSpreadSyncTaskKey   = "GetVideoSpread"
	CreateVideoThumbSyncTaskKey = "CreateVideoThumb"
	CreateVideoM3u8SyncTaskKey  = "CreateVideoM3u8"
)

func init() {
	sync.RegisterProcessor(CreateImageThumbSyncTaskKey, CreateImageThumb)
	sync.RegisterProcessor(GetVideoSpreadSyncTaskKey, GetVideoSpread)
	sync.RegisterProcessor(CreateVideoThumbSyncTaskKey, CreateVideoThumb)
	sync.RegisterProcessor(CreateVideoM3u8SyncTaskKey, CreateVideoM3u8)
}

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
	fileMemeTypes := service.GetMimeTypes(fileType)
	if fileMemeTypes != nil {
		connect.Where("type", fileMemeTypes)
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
	fileModel.Type = uploaded.MimeType
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
	err := sync.NewTask(CreateImageThumbSyncTaskKey, image)
	if err != nil {
		return err
	}
	return nil
}

func CreateImageThumb(param interface{}) {
	imageModel := param.(map[string]interface{})
	openedImage, err := imageLib.Open(fileLib.FullPath(imageModel["path"].(string)))
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	//生成缩略图
	thumb, err := openedImage.CreateThumb(320, 180, "jpg")
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	thumbFile := fileLib.RelativePath(fileLib.Path(thumb.FullPath()))
	//生成预览图
	preview, err := openedImage.CreatePreview(62)
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	previewFile := fileLib.RelativePath(fileLib.Path(preview.FullPath()))
	fileModel := &model.File{}
	connect := fileModel.DB()
	result := connect.Where("id = ?", imageModel["id"]).Updates(map[string]interface{}{"height": openedImage.Height, "width": openedImage.Width, "thumb": thumbFile, "preview": previewFile})
	if result.Error != nil {
		log.Errorf("%err\n", result.Error)
	}
}

func (service File) ProcessVideo(fileModel *model.File) error {
	err := sync.NewTask(GetVideoSpreadSyncTaskKey, fileModel)
	if err != nil {
		return err
	}
	err = sync.NewTask(CreateVideoThumbSyncTaskKey, fileModel)
	if err != nil {
		return err
	}
	err = sync.NewTask(CreateVideoM3u8SyncTaskKey, fileModel)
	if err != nil {
		return err
	}
	return nil
}

func GetVideoSpread(param interface{}) {
	log.Debugf("video spread sync task:%v\n", param)
	videoModel := param.(map[string]interface{})
	openedVideo, err := video.Open(fileLib.FullPath(videoModel["path"].(string)))
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	fileModel := &model.File{
		Width:  openedVideo.Width,
		Height: openedVideo.Height,
	}
	connect := fileModel.DB()
	result := connect.Where("id = ?", videoModel["id"]).Updates(fileModel)
	if result.Error != nil {
		log.Errorf("%err\n", result.Error)
	}
}

func CreateVideoThumb(param interface{}) {
	log.Debugf("video thumb sync task:%v\n", param)
	videoModel := param.(map[string]interface{})
	openedVideo, err := video.Open(fileLib.FullPath(videoModel["path"].(string)))
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	path, progressChannel, err := openedVideo.CreateThumb(320, 180, "jpg", "")
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	videoModelId := int(videoModel["id"].(float64))
	for progress := range progressChannel {
		current := int(progress.GetProgress() * 100)
		err = cache.Set(VideoProgressCacheKey+":"+strconv.Itoa(videoModelId), map[string]interface{}{"action": "createThumb", "current": current, "time": progress.GetCurrentTime()}, ProgressCacheTtl)
		if err != nil {
			log.Debugf("cache thumb progress failed:%err\n", err)
			return
		}
	}
	thumb := fileLib.RelativePath(fileLib.Path(path))
	fileModel := &model.File{}
	connect := fileModel.DB()
	result := connect.Where("id = ?", videoModelId).Update("thumb", thumb)
	if result.Error != nil {
		log.Errorf("生成封面失败：%err\n", result.Error)
	}
}

func CreateVideoM3u8(param interface{}) {
	log.Debugf("video m3u8 sync task:%v\n", param)
	videoModel := param.(map[string]interface{})
	openedVideo, err := video.Open(fileLib.FullPath(videoModel["path"].(string)))
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	path, progressChannel, err := openedVideo.CreateM3u8(720, 576)
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	videoModelId := int(videoModel["id"].(float64))
	for progress := range progressChannel {
		current := int(progress.GetProgress() * 100)
		err = cache.Set(VideoProgressCacheKey+":"+strconv.Itoa(videoModelId), map[string]interface{}{"action": "createM3u8", "current": current, "time": progress.GetCurrentTime()}, ProgressCacheTtl)
		if err != nil {
			log.Debugf("cache m3u8 progress failed:%err\n", err)
			continue
		}
	}
	preview := fileLib.RelativePath(fileLib.Path(path))
	fileModel := &model.File{}
	connect := fileModel.DB()
	result := connect.Where("id = ?", videoModelId).Update("preview", preview)
	if result.Error != nil {
		log.Errorf("生成m3u8失败：%err\n", result.Error)
	}
}
