package model

import (
	"gorm.io/gorm"
	"pmsGo/lib/database"
	"time"
)

type Post struct {
	ID        int       `gorm:"private_key" json:"id"`
	Uuid      string    `gorm:"index" json:"uuid"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	SubTitle  string    `json:"sub_title"`
	Cover     string    `json:"cover"`
	Content   string    `json:"content"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

const (
	PostStatusEnable  = 1
	PostStatusDisable = 2
)

func (model Post) DB() *gorm.DB {
	return database.DB.Model(&model)
}

func (model *Post) Toggle() error {
	if model.Status == PostStatusEnable {
		model.Status = PostStatusDisable
	} else {
		model.Status = PostStatusEnable
	}
	connect := model.DB()
	result := connect.Where("id = ?", model.ID).Save(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
