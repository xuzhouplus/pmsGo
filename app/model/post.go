package model

import (
	"errors"
	"math"
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
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

const (
	PostStatusEnable  = "1"
	PostStatusDisable = "2"
)

var PostModel = &Post{}

func (model Post) List(page interface{}, size interface{}, fields interface{}, like interface{}, enable interface{}, order interface{}) (map[string]interface{}, error) {
	var posts []Post
	connect := database.Query(&Post{})
	if fields != nil {
		connect.Select(fields)
	}
	fetchPage := 1
	if page != nil {
		fetchPage = page.(int)
	}
	if fetchPage < 1 {
		fetchPage = 1
	}
	fetchSize := 10
	if size != nil {
		fetchSize = size.(int)
	}
	connect.Offset((fetchPage - 1) * fetchSize)
	connect.Limit(fetchSize)
	if like != nil {
		connect.Where("title = ?", like)
	}

	if enable != nil {
		connect.Where("status = ?", enable)
	}

	if order != nil {
		for field, sort := range order.(map[string]string) {
			connect.Order(field + " " + sort)
		}
	} else {
		connect.Order("updated_at DESC")
	}
	returnData := make(map[string]interface{})
	result := connect.Find(&posts)
	if result.Error != nil {
		return returnData, errors.New("获取稿件列表失败")
	}
	returnData["posts"] = posts
	returnData["size"] = fetchSize
	returnData["page"] = fetchPage
	var total int64
	result = connect.Count(&total)
	if result.Error != nil {
		return returnData, errors.New("获取稿件总数失败")
	}
	returnData["total"] = total
	returnData["count"] = math.Ceil(float64(int(total) / fetchSize))
	return returnData, nil
}
func (model Post) FindOneById(id int) (*Post, error) {
	one := &Post{}
	connect := database.Query(&Post{})
	connect.Where("id = ?", id)
	connect.Limit(1)
	result := connect.First(&one)
	if result.Error != nil {
		return nil, result.Error
	}
	return one, nil
}
func (model Post) FindOneByUuid(uuid string) (*Post, error) {
	one := &Post{}
	connect := database.Query(&Post{})
	connect.Where("uuid = ?", uuid)
	connect.Limit(1)
	result := connect.First(&one)
	if result.Error != nil {
		return nil, result.Error
	}
	return one, nil
}
func (model Post) Create() {

}
func (model Post) Update() {

}
func (model Post) Delete() {

}
func (model Post) View() {

}
func (model *Post) Toggle() error {
	if model.Status == PostStatusEnable {
		model.Status = PostStatusDisable
	} else {
		model.Status = PostStatusEnable
	}
	connect := database.Query(&Post{})
	result := connect.Save(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
