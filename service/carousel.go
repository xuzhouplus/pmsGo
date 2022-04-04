package service

import (
	"errors"
	"gorm.io/gorm"
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/file/image"
	"pmsGo/lib/log"
	"pmsGo/lib/security/random"
	"pmsGo/lib/sync"
	"pmsGo/model"
	"strconv"
)

const (
	CreateCarouselSyncTaskKey = "CreateCarousel"
)

func init() {
	sync.RegisterProcessor(CreateCarouselSyncTaskKey, CreateCarousel)
}

type Carousel struct {
}

var CarouselService = &Carousel{}

func (service Carousel) List(page int, size int, match map[string]interface{}, titleLike string, order map[string]string) ([]model.Carousel, error) {
	var carousels []model.Carousel
	carouselModel := &model.Carousel{}
	connect := carouselModel.DB()
	if size != 0 {
		connect.Limit(size)
	}
	if page != 0 {
		connect.Offset(page * size)
	}

	if match != nil {
		connect.Where(match)
	}

	if titleLike != "" {
		connect.Where("title like ?", titleLike)
	}
	if order != nil {
		for field, sort := range order {
			connect.Order("`" + field + "` " + sort)
		}
	}
	if connect.Find(&carousels).Error != nil {
		return nil, errors.New("获取轮播图列表失败")
	}
	return carousels, nil
}

func (service Carousel) CreateFiles(uuid string, fileId int) (map[string]interface{}, error) {
	file, err := FileService.FindOne(fileId)
	if err != nil {
		return nil, err
	}
	if uuid == "" {
		uuid = random.Uuid(false)
	}
	err = sync.NewTask(CreateCarouselSyncTaskKey, map[string]interface{}{
		"uuid": uuid,
		"file": fileLib.FullPath(file.Path),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"type":   file.Type,
		"height": 1920,
		"width":  1080,
	}, nil
}

func CreateCarousel(param interface{}) {
	fileModel := param.(map[string]interface{})
	srcImage, err := image.Open(fileModel["file"].(string))
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	uuid := fileModel["uuid"].(string)
	carouselFile, err := srcImage.CreateCarousel(uuid, 1920, 1080, "jpg")
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	thumb, err := carouselFile.CreateThumb(uuid, 320, 180, "")
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	carousel, err := CarouselService.FindByUuid(uuid)
	if err != nil {
		log.Errorf("%err\n", err)
		return
	}
	carousel.Url = fileLib.RelativePath(fileLib.Path(carouselFile.FullPath()))
	carousel.Thumb = fileLib.RelativePath(fileLib.Path(thumb.FullPath()))
	carousel.Width = 1920
	carousel.Height = 1080
	carousel.Status = model.CarouselStatusEnabled
	result := carousel.DB().Save(&carousel)
	if result.Error != nil {
		log.Errorf("%err\n", result.Error)
	}
}

func (service Carousel) Create(fileId int, title string, description string, link string, order int, switchType string, timeout int) (*model.Carousel, error) {
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
	uuid := random.Uuid(false)
	files, err := service.CreateFiles(uuid, fileId)
	if err != nil {
		return nil, err
	}
	carousel.Uuid = uuid
	carousel.Type = files["type"].(string)
	carousel.Height = files["height"].(int)
	carousel.Width = files["width"].(int)
	carousel.FileId = fileId
	carousel.Link = link
	carousel.Title = title
	carousel.Description = description
	carousel.Timeout = timeout
	carousel.Order = order
	if switchType == "" {
		switchType = model.SwitchTypeWebgl
	}
	carousel.SwitchType = switchType
	carousel.Status = model.CarouselStatusPreparing
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

func (service Carousel) Update(id int, fileId int, title string, description string, link string, order int, switchType string, timeout int) (*model.Carousel, error) {
	carousel := &model.Carousel{}
	db := carousel.DB()
	db.Where("id = ?", id)
	db.Limit(1)
	result := db.First(&carousel)
	if result.Error != nil {
		return nil, result.Error
	}
	db.Begin()
	if carousel.FileId != fileId {
		files, err := service.CreateFiles(carousel.Uuid, fileId)
		if err != nil {
			return nil, err
		}
		carousel.Type = files["type"].(string)
		carousel.Height = files["height"].(int)
		carousel.Width = files["width"].(int)
		carousel.Status = model.CarouselStatusPreparing
	}
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
	if switchType == "" {
		switchType = model.SwitchTypeWebgl
	}
	carousel.SwitchType = switchType
	if timeout < 3 {
		timeout = 3
	}
	if timeout > 10 {
		timeout = 10
	}
	carousel.Timeout = timeout
	connect := carousel.DB()
	result = connect.Save(&carousel)
	if result.Error != nil {
		db.Rollback()
		return nil, result.Error
	}
	db.Commit()
	return carousel, nil
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
