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
	Name        string `json:"name"`
	Thumb       string `json:"thumb"`
	Path        string `json:"path"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Description string `json:"description"`
	Preview     string `json:"preview"`
}

func (model File) DB() *gorm.DB {
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
