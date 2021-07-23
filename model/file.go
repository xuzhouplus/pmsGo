package model

import (
	"gorm.io/gorm"
	"pmsGo/lib/database"
	"pmsGo/lib/image"
)

type File struct {
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

func (model File) DB() *gorm.DB {
	return database.DB.Model(&model)
}
func (model File) RemoveFile() error {
	return image.Remove(image.FullPath(model.Path))
}
func (model File) RemoveThumb() error {
	return image.Remove(image.FullPath(model.Thumb))
}
func (model File) RemovePreview() error {
	return image.Remove(image.FullPath(model.Preview))
}
