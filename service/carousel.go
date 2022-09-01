package service

import (
	"errors"
	"gorm.io/gorm"
	fileLib "pmsGo/lib/file"
	"pmsGo/lib/security/json"
	"pmsGo/lib/security/random"
	"pmsGo/lib/sync"
	"pmsGo/model"
	"pmsGo/worker"
	"strconv"
)

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
	} else {
		connect.Order("`order` asc ")
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
	task, err := sync.NewTask(worker.CarouselWorkerName, map[string]interface{}{
		"uuid":   uuid,
		"path":   fileLib.FullPath(file.Path),
		"width":  1920,
		"height": 1080,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"task_id": task.UUID,
		"type":    file.Type,
		"height":  1920,
		"width":   1080,
	}, nil
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

func (service Carousel) SetTitleStyle(id int, style *model.CaptionStyle) error {
	carouselRecord, err := service.FindById(id)
	if err != nil {
		return err
	}
	carouselRecord.TitleStyle, _ = json.Encode(style)
	return carouselRecord.Save()
}

func (service Carousel) SetDescriptionStyle(id int, style *model.CaptionStyle) error {
	carouselRecord, err := service.FindById(id)
	if err != nil {
		return err
	}
	carouselRecord.DescriptionStyle, _ = json.Encode(style)
	return carouselRecord.Save()
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
		result = connect.Where("`order` > ?", from).Where("`order` <= ?", to).Update("order", gorm.Expr("`order` - ?", 1))
	} else {
		result = connect.Where("`order` >= ?", to).Where("`order` < ?", from).Update("order", gorm.Expr("`order` + ?", 1))
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (service Carousel) SortOrder(sortedOrder map[int]int) error {
	carousel := &model.Carousel{}
	connect := carousel.DB()
	connect.Begin()
	for id, order := range sortedOrder {
		result := carousel.DB().Where("`id` = ?", id).Update("order", order)
		if result.Error != nil {
			connect.Rollback()
			return result.Error
		}
	}
	connect.Commit()
	return nil
}

func (service Carousel) UpdateCaptionStyle(id interface{}, title interface{}, link interface{}, titleStyle interface{}, description interface{}, descriptionStyle interface{}, switchType interface{}) (*model.Carousel, error) {
	carousel, err := service.FindById(int(id.(float64)))
	if err != nil {
		return nil, err
	}
	carousel.Title = title.(string)
	carousel.Link = link.(string)
	carousel.TitleStyle, _ = json.Encode(titleStyle)
	carousel.Description = description.(string)
	carousel.DescriptionStyle, _ = json.Encode(descriptionStyle)
	carousel.SwitchType = switchType.(string)
	err = carousel.Save()
	if err != nil {
		return nil, err
	}
	return carousel, nil
}

func (service Carousel) SetSwitchType(id interface{}, switchType interface{}) (*model.Carousel, error) {
	carousel, err := service.FindById(int(id.(float64)))
	if err != nil {
		return nil, err
	}
	carousel.SwitchType = switchType.(string)
	err = carousel.Save()
	if err != nil {
		return nil, err
	}
	return carousel, nil
}
