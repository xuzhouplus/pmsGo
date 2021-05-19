package service

import (
	"errors"
	"pmsGo/app/model"
	"pmsGo/lib/database"
	"pmsGo/lib/helper/image"
	"pmsGo/lib/security/random"
	"strconv"
)

type Carousel struct {
}

var CarouselService = &Carousel{}

func (service Carousel) List(page interface{}, size interface{}, fields interface{}, like interface{}, order interface{}) ([]model.Carousel, error) {
	var carousels []model.Carousel

	connect := database.Query(&model.Carousel{})

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
func (service Carousel) Create(fileId int, title string, description string, link string, order int) (*model.Carousel, error) {
	carouselLimit, _ := strconv.Atoi(SettingService.GetSetting(model.SettingKeyCarouselLimit))
	connect := database.Query(&model.Carousel{})
	var count int64
	result := connect.Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}
	if carouselLimit == int(count) {
		return nil, errors.New("the carousels number is reached limit")
	}
	one, err := FileService.FindOne(fileId)
	if err != nil {
		return nil, err
	}
	carousel := &model.Carousel{}
	carousel.Type = one.Type
	carousel.Link = link
	carousel.Title = title
	carousel.Description = description
	carousel.Order = order
	file, err := image.Open(image.FullPath(one.Path))
	if err != nil {
		return nil, err
	}
	carouselFile, err := file.CreateCarousel(1920, 1080, "jpg")
	if err != nil {
		return nil, err
	}
	carousel.Url = image.RelativePath(image.Path(carouselFile.FullPath()))
	thumb, err := carouselFile.CreateThumb(320, 180, "")
	if err != nil {
		return nil, err
	}
	carousel.Thumb = image.RelativePath(image.Path(thumb.FullPath()))
	carousel.Height = carouselFile.Height
	carousel.Width = carouselFile.Width
	carousel.FileId = fileId
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
func (service Carousel) Preview(fileId int) (string, error) {
	one, err := FileService.FindOne(fileId)
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

func (service Carousel) Delete(id int) error {
	carousel := &model.Carousel{}
	connect := database.Query(&model.Carousel{})
	connect.Where("id = ?", id)
	connect.Limit(1)
	result := connect.Find(&carousel)
	if result.Error != nil {
		return result.Error
	}
	err := image.Remove(image.FullPath(carousel.Thumb))
	if err != nil {
		return err
	}
	result = connect.Delete(&carousel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
