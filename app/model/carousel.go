package model

import (
	"errors"
	"pmsGo/lib/database"
	"pmsGo/lib/helper/image"
	"pmsGo/lib/security"
	"strconv"
)

type Carousel struct {
	ID          int    `gorm:"private_key" json:"id"`
	Uuid        string `gorm:"index" json:"uuid"`
	FileId      int    `json:"field_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Thumb       string `json:"thumb"`
	Link        string `json:"link"`
}

var CarouselModel = &Carousel{}

func (model Carousel) List(page interface{}, size interface{}, fields interface{}, like interface{}, order interface{}) ([]Carousel, error) {
	var carousels []Carousel

	connect := database.Query(&Carousel{})

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
			connect.Order("`" + field + "` " + sort)
		}
	}

	if connect.Find(&carousels).Error != nil {
		return carousels, errors.New("获取轮播图列表失败")
	}
	for i, carousel := range carousels {
		carousel.Url = image.FullUrl(carousel.Url)
		carousel.Thumb = image.FullUrl(carousel.Thumb)
		carousels[i] = carousel
	}
	return carousels, nil
}

func (model *Carousel) Create(fileId int, title string, description string, link string, order int) error {
	carouselLimit, _ := strconv.Atoi(CarouselSettingModel.GetSetting(SettingKeyCarouselLimit))
	connect := database.Query(&Carousel{})
	var count int64
	result := connect.Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if carouselLimit == int(count) {
		return errors.New("the carousels number is reached limit")
	}
	one, err := FileModel.FindOne(fileId)
	if err != nil {
		return err
	}
	model.Type = one.Type
	model.Link = link
	model.Title = title
	model.Description = description
	model.Order = order
	file, err := image.Open(image.FullPath(one.Path))
	if err != nil {
		return err
	}
	carousel, err := file.CreateCarousel(1920, 1080, "jpg")
	if err != nil {
		return err
	}
	model.Url = image.RelativePath(image.Path(carousel.FullPath()))
	thumb, err := carousel.CreateThumb(320, 180, "")
	if err != nil {
		return err
	}
	model.Thumb = image.RelativePath(image.Path(thumb.FullPath()))
	model.Height = carousel.Height
	model.Width = carousel.Width
	model.FileId = fileId
	model.Uuid = security.Uuid(false)
	result = connect.Create(&model)
	if result.Error != nil {
		return result.Error
	}
	err = model.SetOrder(order)
	if err != nil {
		return err
	}
	return nil
}

func (model *Carousel) SetOrder(order interface{}) error {
	var carousels []Carousel
	connect := database.Query(&Carousel{})
	connect.Begin()
	if order == nil {
		connect.Where("order > ?", model.Order)
		result := connect.Find(&carousels)
		if result.Error != nil {
			return result.Error
		}
		for _, carousel := range carousels {
			carousel.Order = carousel.Order - 1
			result = connect.Save(&carousel)
			if result.Error != nil {
				connect.Rollback()
				return result.Error
			}
		}
		connect.Commit()
		return nil
	}
	orderInt := order.(int)
	if orderInt < model.Order {
		connect.Where("order < ?", model.Order).Where("order > ?", orderInt-1)
		result := connect.Find(&carousels)
		if result.Error != nil {
			return result.Error
		}
		for _, carousel := range carousels {
			result = connect.Where("id", carousel.ID).Update("order", carousel.Order+1)
			if result.Error != nil {
				connect.Rollback()
				return result.Error
			}
		}
	} else if orderInt > model.Order {
		connect.Where("order < ?", orderInt).Where("order > ?", model.Order)
		result := connect.Find(&carousels)
		if result.Error != nil {
			return result.Error
		}
		for _, carousel := range carousels {
			result = connect.Where("id", carousel.ID).Update("order", carousel.Order-1)
			result = connect.Save(&carousel)
			if result.Error != nil {
				connect.Rollback()
				return result.Error
			}
		}
	} else {
		connect.Order("`order` asc")
		result := connect.Find(&carousels)
		if result.Error != nil {
			return result.Error
		}
		orderIndex := 1
		for _, carousel := range carousels {
			result = connect.Where("id", carousel.ID).Update("order", orderIndex)
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

func (model Carousel) Preview(fileId int) (string, error) {
	one, err := FileModel.FindOne(fileId)
	if err != nil {
		return "", err
	}
	open, err := image.Open(image.FullPath(one.Preview))
	if err != nil {
		return "", err
	}
	carousel, err := open.CreateCarousel(1920, 1080, "jpg")
	if err != nil {
		return "", err
	}
	return image.FullUrl(image.RelativePath(image.Path(carousel.FullPath()))), nil
}

func (model *Carousel) Delete(id int) error {
	connect := database.Query(&Carousel{})
	connect.Where("id = ?", id)
	connect.Limit(1)
	result := connect.Find(&model)
	if result.Error != nil {
		return result.Error
	}
	err := image.Remove(image.FullPath(model.Thumb))
	if err != nil {
		return err
	}
	result = connect.Delete(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
