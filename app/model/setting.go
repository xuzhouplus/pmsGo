package model

import (
	"pmsGo/lib/database"
)

const SettingTypeIsPrivate = 1

type Setting struct {
	database.Model
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

