package service

import (
	"errors"
	"github.com/floostack/transcoder"
	"math"
	"pmsGo/lib/cache"
	fileLib "pmsGo/lib/file"
	imageLib "pmsGo/lib/file/image"
	"pmsGo/lib/file/video"
	"pmsGo/lib/log"
	"pmsGo/lib/sync"
	"pmsGo/model"
	"strconv"
)

const (
	VideoProgressCacheKey = "progress:video"
	ImageProgressCacheKey = "progress:image"
	ProgressCacheTtl      = 60 * 60 * 24
)

type File struct {
}

var FileService = &File{}

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
		connect.Where("type = ?", fileType)
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
	result := connect.Delete(&one)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (service File) ProcessImage(image *model.File) {
	syncTask := sync.NewTask(image, func(uuid string, param interface{}) (string, error) {
		imageModel := param.(*model.File)
		openedImage, err := imageLib.Open(fileLib.FullPath(imageModel.Path))
		if err != nil {
			return uuid, err
		}
		//文件幅面大小
		imageModel.Height = openedImage.Height
		imageModel.Width = openedImage.Width
		//生成缩略图
		thumb, err := openedImage.CreateThumb(320, 180, "jpg")
		if err != nil {
			return uuid, err
		}
		imageModel.Thumb = fileLib.RelativePath(fileLib.Path(thumb.FullPath()))
		//生成预览图
		preview, err := openedImage.CreatePreview(62)
		if err != nil {
			return uuid, err
		}
		imageModel.Preview = fileLib.RelativePath(fileLib.Path(preview.FullPath()))
		connect := imageModel.DB()
		result := connect.Select("Height", "Width", "Thumb", "Preview").Updates(imageModel)
		if result.Error != nil {
			return uuid, result.Error
		}
		return uuid, nil
	})
	sync.AddTask(syncTask)
}

func (service File) ProcessVideo(image *model.File) {
	syncTask := sync.NewTask(image, func(uuid string, param interface{}) (string, error) {
		log.Debugf("sync task:%v\n", uuid)
		videoModel := param.(*model.File)
		log.Debugf("sync task:%v\n", fileLib.FullPath(videoModel.Path))
		openedVideo, err := video.Open(fileLib.FullPath(videoModel.Path))
		if err != nil {
			return uuid, err
		}
		//文件幅面大小
		videoModel.Height = openedVideo.Height
		videoModel.Width = openedVideo.Width
		thumbFile, err := openedVideo.CreateThumb(320, 180, "jpg", "", func(progresses <-chan transcoder.Progress) {
			go func() {
				for progress := range progresses {
					time, err := strconv.Atoi(progress.GetCurrentTime())
					if err != nil {
						time = 0
					}
					current := int(progress.GetProgress())
					err = cache.Set(VideoProgressCacheKey+":"+strconv.Itoa(videoModel.ID), map[string]interface{}{"action": "createThumb", "current": current, "time": time}, ProgressCacheTtl)
					if err != nil {
						return
					}
					if current==100{
						thumb, err := imageLib.Open(thumbFile)
						if err != nil {
							return uuid, err
						}
						videoModel.Thumb = fileLib.RelativePath(fileLib.Path(thumb.FullPath()))
					}
				}
			}()
		})
		if err != nil {
			return uuid, err
		}

		m3u8, err := openedVideo.CreateM3u8(720, 576, func(progresses <-chan transcoder.Progress) {
			go func() {
				for progress := range progresses {
					time, err := strconv.Atoi(progress.GetCurrentTime())
					if err != nil {
						time = 0
					}
					current := int(progress.GetProgress())
					err = cache.Set(VideoProgressCacheKey+":"+strconv.Itoa(videoModel.ID), map[string]interface{}{"action": "createM3u8", "current": current, "time": time}, ProgressCacheTtl)
					if err != nil {
						return
					}
					if current==100{
						
					}
				}
			}()
		})
		if err != nil {
			return "", err
		}
		videoModel.Preview = fileLib.RelativePath(fileLib.Path(m3u8))
		connect := videoModel.DB()
		result := connect.Select("Height", "Width", "Thumb", "Preview").Updates(videoModel)
		if result.Error != nil {
			return uuid, result.Error
		}
		return uuid, nil
	})
	sync.AddTask(syncTask)
}
