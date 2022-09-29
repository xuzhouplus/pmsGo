package service

import (
	"errors"
	"gorm.io/gorm"
	"math"
	"pmsGo/lib/security/random"
	model "pmsGo/model"
)

type Post struct {
}

var PostService = &Post{}

func (service Post) List(page int, size int, fields []string, like string, enable int, order map[string]string) (map[string]interface{}, error) {
	var posts []model.Post
	postModel := &model.Post{}
	connect := postModel.DB()
	if fields != nil {
		connect.Select(fields)
	}
	if page < 0 {
		page = 0
	}
	if size == 0 {
		size = 10
	}
	connect.Offset((page - 1) * size)
	connect.Limit(size)
	if like != "" {
		connect.Where("title = ?", like)
	}
	if enable != 0 {
		connect.Where("status = ?", enable)
	}
	if order != nil {
		for field, sort := range order {
			connect.Order("`" + field + "` " + sort)
		}
	} else {
		connect.Order("`updated_at` DESC")
	}
	returnData := make(map[string]interface{})
	result := connect.Find(&posts)
	if result.Error != nil {
		return returnData, errors.New("获取稿件列表失败")
	}
	returnData["posts"] = posts
	returnData["size"] = size
	returnData["page"] = page
	var total int64
	connect.Offset(-1)
	connect.Limit(-1)
	result = connect.Count(&total)
	if result.Error != nil {
		return returnData, errors.New("获取稿件总数失败")
	}
	returnData["total"] = total
	returnData["count"] = math.Ceil(float64(total) / float64(size))
	return returnData, nil
}
func (service Post) FindOneById(id int) (*model.Post, error) {
	one := &model.Post{}
	connect := one.DB()
	connect.Where("id = ?", id)
	connect.Limit(1)
	result := connect.First(&one)
	if result.Error != nil {
		return nil, result.Error
	}
	return one, nil
}
func (service Post) FindOneByUuid(uuid string) (*model.Post, error) {
	one := &model.Post{}
	connect := one.DB()
	connect.Where("uuid = ?", uuid)
	connect.Limit(1)
	result := connect.First(&one)
	if result.Error != nil {
		return nil, result.Error
	}
	return one, nil
}
func (service Post) Save(postData map[string]interface{}) (*model.Post, error) {
	var err error
	data := &model.Post{}
	if postData["id"] != nil {
		data, err = service.FindOneById(int(postData["id"].(float64)))
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("post not exists")
		}
	}
	data.Type = postData["type"].(string)
	data.Title = postData["title"].(string)
	data.SubTitle = postData["sub_title"].(string)
	data.Content = postData["content"].(string)
	data.Cover = postData["cover"].(string)
	data.Status = int(postData["status"].(float64))
	newRecord := false
	if data.Uuid == "" {
		newRecord = true
		data.Uuid = random.Uuid(false)
	}
	connect := data.DB()
	var result *gorm.DB
	if newRecord {
		result = connect.Create(&data)
	} else {
		result = connect.Where("id = ?", data.ID).Save(&data)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}
func (service Post) Delete(uuid string) error {
	post, _ := service.FindOneByUuid(uuid)
	if post == nil {
		return errors.New("post not exists")
	}
	connect := post.DB()
	result := connect.Where("id = ?", post.ID).Delete(&post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
