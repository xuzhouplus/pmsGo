package model

import (
	"gorm.io/gorm"
	"pmsGo/lib/database"
)

const SettingTypeIsPrivate = 1

type Setting struct {
	ID          int    `gorm:"private_key" json:"id"`
	Key         string `gorm:"uniqueIndex" json:"key"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Private     int    `json:"private"`
	Value       string `json:"value"`
	Options     string `json:"options"`
	Description string `json:"description"`
	Required    int    `json:"required"`
}

func (model *Setting) DB() *gorm.DB {
	return database.DB.Model(&model)
}
