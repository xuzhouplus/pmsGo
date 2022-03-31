package model

import (
	"gorm.io/gorm"
	"pmsGo/lib/database"
)

const (
	SwitchTypeWebgl    = "webgl"
	SwitchTypeSeparate = "separate"
	SwitchTypeSlide    = "slide"
)

type Carousel struct {
	ID          int    `gorm:"private_key" json:"id"`
	Uuid        string `gorm:"index" json:"uuid"`
	FileId      int    `json:"file_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Thumb       string `json:"thumb"`
	Link        string `json:"link"`
	SwitchType  string `json:"switch_type"`
	Timeout     int    `json:"timeout"`
	Duration    int    `json:"duration"`
}

func (model *Carousel) DB() *gorm.DB {
	return database.DB.Model(&model)
}

func (model *Carousel) SetOrder(order interface{}) error {
	var carousels []Carousel
	connect := model.DB()
	if order == nil {
		result := connect.Where("order >= ?", model.Order).Update("order", gorm.Expr("order + ?", 1))
		return result.Error
	}
	connect.Begin()
	orderInt := order.(int)
	if orderInt < model.Order {
		result := connect.Model(&Carousel{}).Where("order < ?", model.Order).Where("order > ?", orderInt-1).Find(&carousels)
		if result.Error != nil {
			return result.Error
		}
		for _, carousel := range carousels {
			result = connect.Model(&Carousel{}).Where("id", carousel.ID).Update("order", carousel.Order+1)
			if result.Error != nil {
				connect.Rollback()
				return result.Error
			}
		}
	} else if orderInt > model.Order {
		result := connect.Where("order < ?", orderInt).Where("order > ?", model.Order).Find(&carousels)
		if result.Error != nil {
			return result.Error
		}
		for _, carousel := range carousels {
			result = connect.Model(&Carousel{}).Where("id", carousel.ID).Update("order", carousel.Order-1).Save(&carousel)
			if result.Error != nil {
				connect.Rollback()
				return result.Error
			}
		}
	} else {
		result := connect.Model(&Carousel{}).Order("`order` asc").Find(&carousels)
		if result.Error != nil {
			return result.Error
		}
		orderIndex := 1
		for _, carousel := range carousels {
			result = connect.Model(&Carousel{}).Where("id", carousel.ID).Update("order", orderIndex)
			if result.Error != nil {
				connect.Rollback()
				return result.Error
			}
			orderIndex++
		}
	}
	model.Order = orderInt
	connect.Save(&model)
	connect.Commit()
	return nil
}
