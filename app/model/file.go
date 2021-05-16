package model

import (
	"errors"
	"math"
	"pmsGo/lib/database"
	"pmsGo/lib/helper/image"
)

type file struct {
	ID          int    `gorm:"private_key" json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Thumb       string `json:"thumb"`
	Path        string `json:"path"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Description string `json:"description"`
	Preview     string `json:"preview"`
}

var FileModel = &file{}

func (model file) List(page interface{}, limit interface{}, fields interface{}, fileType interface{}, name interface{}) (map[string]interface{}, error) {
	var files []file
	connect := database.Query(&file{})
	fetchPage := 1
	if page != nil {
		fetchPage = page.(int)
	}
	if fetchPage < 1 {
		fetchPage = 1
	}
	fetchSize := 10
	if limit != nil {
		fetchSize = limit.(int)
	}
	connect.Offset((fetchPage - 1) * fetchSize)
	connect.Limit(fetchSize)
	if fields != nil {
		connect.Select(fields)
	}
	if name != nil {
		connect.Where("name like ?", name)
	}
	if fileType != nil {
		connect.Where("type = ?", fileType.(string))
	}
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
	returnData["size"] = fetchSize
	returnData["page"] = fetchPage
	var total int64
	connect.Count(&total)
	returnData["total"] = total
	returnData["count"] = math.Ceil(float64(int(total) / fetchSize))
	return returnData, nil
}
func (model *file) Upload(uploaded *image.Upload, name string, description string) error {
	model.Name = name
	model.Description = description
	model.Path = image.RelativePath(uploaded.Path())
	model.Type = uploaded.MimeType
	filePath := uploaded.Path()
	fileImage, err := image.Open(string(filePath))
	if err != nil {
		return err
	}
	model.Width = fileImage.Width
	model.Height = fileImage.Height
	thumb, err := fileImage.CreateThumb(320, 180, "jpg")
	if err != nil {
		return err
	}
	model.Thumb = image.RelativePath(image.Path(thumb.FullPath()))
	preview, err := fileImage.CreatePreview(62)
	if err != nil {
		return err
	}
	model.Preview = image.RelativePath(image.Path(preview.FullPath()))
	connect := database.Query(&file{})
	result := connect.Create(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (model file) FindOne(id int) (*file, error) {
	one := &file{}
	connect := database.Query(&file{})
	connect.Where("id = ?", id)
	connect.Limit(1)
	err := connect.Find(&one).Error
	if err != nil {
		return nil, err
	}
	return one, nil
}
func (model file) RemoveFile() error {
	return image.Remove(image.FullPath(model.Path))
}
func (model file) RemoveThumb() error {
	return image.Remove(image.FullPath(model.Thumb))
}
func (model file) RemovePreview() error {
	return image.Remove(image.FullPath(model.Preview))
}
func (model file) Delete(id int) error {
	var one file
	connect := database.Query(&file{})
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
