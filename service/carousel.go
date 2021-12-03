package service

import (
	"errors"
	"gorm.io/gorm"
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/file/image"
	"pmsGo/lib/security/random"
	"pmsGo/model"
	"strconv"
)

type Carousel struct {
}

var CarouselService = &Carousel{}

func (service Carousel) List(page interface{}, size interface{}, fields interface{}, like interface{}, order interface{}) ([]model.Carousel, error) {
	var carousels []model.Carousel
	carouselModel := &model.Carousel{}
	connect := carouselModel.DB()
	if size != nil {
		connect.Limit(size.(int))
	}
	if page != nil {
		connect.Offset(page.(int) * size.(int))
	}
	if fields != nil {
		connect.Select(fields)
	}
	if like != nil && like != "" {
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
		carousel.Url = fileLib.FullUrl(carousel.Url)
		carousel.Thumb = fileLib.FullUrl(carousel.Thumb)
		carousels[i] = carousel
	}
	return carousels, nil
}

func (service Carousel) CreateFiles(fileId int) (map[string]interface{}, error) {
	file, err := FileService.FindOne(fileId)
	if err != nil {
		return nil, err
	}
	srcImage, err := image.Open(fileLib.FullPath(file.Path))
	if err != nil {
		return nil, err
	}
	carouselFile, err := srcImage.CreateCarousel(1920, 1080, "jpg")
	if err != nil {
		return nil, err
	}
	thumb, err := carouselFile.CreateThumb(320, 180, "")
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"type":   file.Type,
		"url":    fileLib.RelativePath(fileLib.Path(carouselFile.FullPath())),
		"thumb":  fileLib.RelativePath(fileLib.Path(thumb.FullPath())),
		"height": carouselFile.Height,
		"width":  carouselFile.Width,
	}, nil
}

func (service Carousel) Create(fileId int, title string, description string, link string, order int) (*model.Carousel, error) {
	carouselLimit, _ := strconv.Atoi(SettingService.GetSetting(model.SettingKeyCarouselLimit))
	carouselModel := &model.Carousel{}
	connect := carouselModel.DB()
	var count int64
	result := connect.Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}
	if carouselLimit == int(count) {
		return nil, errors.New("the carousels number is reached limit")
	}
	carousel := &model.Carousel{}
	files, err := service.CreateFiles(fileId)
	if err != nil {
		return nil, err
	}
	carousel.Type = files["type"].(string)
	carousel.Url = files["url"].(string)
	carousel.Thumb = files["thumb"].(string)
	carousel.Height = files["height"].(int)
	carousel.Width = files["width"].(int)
	carousel.FileId = fileId
	carousel.Link = link
	carousel.Title = title
	carousel.Description = description
	carousel.Order = order
	carousel.Uuid = random.Uuid(false)
	result = connect.Create(&carousel)
	if result.Error != nil {
		return nil, result.Error
	}
	err = carousel.SetOrder(order)
	if err != nil {
		return nil, err
	}
	return carousel, nil
}

func (service Carousel) Update(id int, fileId int, title string, description string, link string, order int) (*model.Carousel, error) {
	carousel := &model.Carousel{}
	db := carousel.DB()
	db.Where("id = ?", id)
	db.Limit(1)
	result := db.First(&carousel)
	if result.Error != nil {
		return nil, result.Error
	}
	db.Begin()
	files, err := service.CreateFiles(fileId)
	if err != nil {
		return nil, err
	}
	carousel.Type = files["type"].(string)
	carousel.Url = files["url"].(string)
	carousel.Thumb = files["thumb"].(string)
	carousel.Height = files["height"].(int)
	carousel.Width = files["width"].(int)
	carousel.FileId = fileId
	carousel.Link = link
	if title != "" {
		carousel.Title = title
	}
	if description != "" {
		carousel.Description = description
	}
	if order != 0 {
		if carousel.Order != order {
			err := service.UpdateOrder(carousel.Order, order)
			if err != nil {
				return nil, err
			}
		}
		carousel.Order = order

	}
	connect := carousel.DB()
	result = connect.Save(&carousel)
	if result.Error != nil {
		db.Rollback()
		return nil, result.Error
	}
	db.Commit()
	return carousel, nil
}

func (service Carousel) Preview(fileId int) (string, error) {
	one, err := FileService.FindOne(fileId)
	if err != nil {
		return "", err
	}
	open, err := image.Open(fileLib.FullPath(one.Preview))
	if err != nil {
		return "", err
	}
	carousel, err := open.CreateCarousel(1920, 1080, "jpg")
	if err != nil {
		return "", err
	}
	return fileLib.FullUrl(fileLib.RelativePath(fileLib.Path(carousel.FullPath()))), nil
}

func (service Carousel) Delete(id int) error {
	carousel := &model.Carousel{}
	connect := carousel.DB()
	connect.Where("id = ?", id)
	connect.Limit(1)
	result := connect.Find(&carousel)
	if result.Error != nil {
		return result.Error
	}
	err := fileLib.Remove(fileLib.FullPath(carousel.Url))
	if err != nil {
		return err
	}
	err = fileLib.Remove(fileLib.FullPath(carousel.Thumb))
	if err != nil {
		return err
	}
	result = connect.Delete(&carousel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (service Carousel) FindById(id int) (*model.Carousel, error) {
	one := &model.Carousel{}
	connect := one.DB()
	connect.Where("id = ?", id)
	connect.Limit(1)
	result := connect.First(&one)
	if result.Error != nil {
		return nil, result.Error
	}
	return one, nil
}

func (service Carousel) FindByUuid(uuid string) (*model.Carousel, error) {
	one := &model.Carousel{}
	connect := one.DB()
	connect.Where("uuid = ?", uuid)
	connect.Limit(1)
	result := connect.First(&one)
	if result.Error != nil {
		return nil, result.Error
	}
	return one, nil
}

func (service Carousel) FindByOrder(order int) (*model.Carousel, error) {
	one := &model.Carousel{}
	connect := one.DB()
	connect.Where("order = ?", order)
	connect.Limit(1)
	result := connect.First(&one)
	if result.Error != nil {
		return nil, result.Error
	}
	return one, nil
}

func (service Carousel) UpdateOrder(from int, to int) error {
	one := &model.Carousel{}
	connect := one.DB()
	var result *gorm.DB
	if from < to {
		result = connect.Where("order > ?", from).Where("order <= ?", to).Update("order", gorm.Expr("order - ?", 1))
	} else {
		result = connect.Where("order >= ?", to).Where("order < ?", from).Update("order", gorm.Expr("order + ?", 1))
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
