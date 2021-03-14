package model

import (
	"errors"
	"pmsGo/lib/database"
)

type carousel struct {
	ID          int    `gorm:"private_key" json:"id"`
	Uuid        string `gorm:"index" json:"uuid"`
	FileId      int `json:"field_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Width       int `json:"width"`
	Height      int `json:"height"`
	Description string `json:"description"`
	Order       int `json:"order"`
	Thumb       string `json:"thumb"`
	Link        string `json:"link"`
}

var Carousel = &carousel{}

func (model carousel) List(page interface{}, size interface{}, fields interface{}, like interface{}, order interface{}) ([]carousel, error) {
	var carousels []carousel

	connect := database.Connect(&carousel{})

	if size != nil {
		connect.Limit(size.(int))
	}

	if page != nil {
		connect.Offset(page.(int) * size.(int))
	}

	if fields != nil {
		connect.Select(fields)
	}
	if like != nil {
		connect.Where("title like ?", like.(string))
	}

	if order != nil {
		for field, sort := range order.(map[string]string) {
			connect.Order(field + " " + sort)
		}
	}

	if connect.Find(&carousels).Error != nil {
		return carousels, errors.New("获取轮播图列表失败")
	}

	return carousels, nil
}
