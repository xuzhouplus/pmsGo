package model

import (
	"errors"
	"fmt"
	"pmsGo/lib/database"
	"time"
)

type post struct {
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

var Post = &post{}

func (model post) List(page interface{}, size interface{}, fields interface{}, like interface{}, enable interface{}, order interface{}) (map[string]interface{}, error) {
	var posts []post
	connect := database.Connect(&post{})
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
	fmt.Println(result.RowsAffected)
	returnData["posts"] = posts
	returnData["size"] = fetchSize
	returnData["page"] = fetchPage
	//returnData["count"] = connect.Count(&posts)
	return returnData, nil
}
