package model

import (
	"gorm.io/gorm"
	"pmsGo/lib/config"
	"pmsGo/lib/database"
	fileLib "pmsGo/lib/file"
)

type File struct {
	ID          int    `gorm:"private_key" json:"id"`
	Uuid        string `json:"uuid"`
	Type        string `json:"type"`
	Extension   string `json:"extension"`
	Name        string `json:"name"`
	Poster      string `json:"poster"`
	Thumb       string `json:"thumb"`
	Path        string `json:"path"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Description string `json:"description"`
	Preview     string `json:"preview"`
	Status      int    `json:"status"`
}

const (
	StatusUploaded   = 0
	StatusProcessing = 1
	StatusEnabled    = 2
)

func (model *File) DB() *gorm.DB {
	return database.DB.Model(&model)
}

func (model File) RemoveFile() error {
	return fileLib.Remove(fileLib.FullPath(model.Path))
}

func (model File) RemoveThumb() error {
	thumbFile := fileLib.FullPath(model.Thumb)
	if thumbFile == "" {
		return nil
	}
	return fileLib.Remove(thumbFile)
}

func (model File) RemovePoster() error {
	if model.Poster == "" {
		return nil
	}
	posterFile := fileLib.FullPath(model.Poster)
	if posterFile == "" {
		return nil
	}
	return fileLib.Remove(posterFile)
}

func (model File) RemovePreview() error {
	previewFile := fileLib.FullPath(model.Preview)
	if previewFile == "" {
		return nil
	}
	return fileLib.Remove(previewFile)
}

func (model File) RemoveDir() error {
	return fileLib.Remove(config.Config.Web.Upload.Path + "/" + model.Uuid)
}

func (model *File) Update(field string, value interface{}) error {
	db := model.DB()
	result := db.Update(field, value)
	return result.Error
}
