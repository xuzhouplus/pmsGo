package service

import (
	"errors"
	"math"
	"pmsGo/app/model"
	"pmsGo/lib/helper/image"
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
		connect.Where("name like ?", name)
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
			f.Path = image.FullUrl(f.Path)
			f.Thumb = image.FullUrl(f.Thumb)
			f.Preview = image.FullUrl(f.Preview)
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
func (service File) Upload(uploaded *image.Instance, name string, description string) (*model.File, error) {
	file := &model.File{}
	file.Name = name
	file.Description = description
	file.Path = image.RelativePath(uploaded.Path())
	file.Type = uploaded.MimeType
	filePath := uploaded.Path()
	fileImage, err := image.Open(string(filePath))
	if err != nil {
		return nil, err
	}
	file.Width = fileImage.Width
	file.Height = fileImage.Height
	channel := make(chan map[string]string, 2)
	go func() {
		thumb, err := fileImage.CreateThumb(320, 180, "jpg")
		if err != nil {
			channel <- map[string]string{"error": err.Error()}
		} else {
			thumb := image.RelativePath(image.Path(thumb.FullPath()))
			channel <- map[string]string{"thumb": thumb}
		}
	}()
	go func() {
		preview, err := fileImage.CreatePreview(62)
		if err != nil {
			channel <- map[string]string{"error": err.Error()}
		}
		prev := image.RelativePath(image.Path(preview.FullPath()))
		channel <- map[string]string{"preview": prev}
	}()
	for i := range channel {
		if i["error"] != "" {
			close(channel)
			return nil, errors.New(i["error"])
		}
		if i["preview"] != "" {
			file.Preview = i["preview"]
		} else if i["thumb"] != "" {
			file.Thumb = i["thumb"]
		}
		if file.Preview != "" && file.Thumb != "" {
			close(channel)
		}
	}
	connect := file.DB()
	result := connect.Create(&file)
	if result.Error != nil {
		return nil, result.Error
	}
	return file, nil
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
