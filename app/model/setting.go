package model

import (
	"errors"
	"pmsGo/lib/database"
)

const TypeIsPrivate = 1

type setting struct {
	database.Model
	ID          int    `gorm:"private_key"`
	Key         string `gorm:"uniqueIndex"`
	Name        string
	Type        string
	Private     int
	Value       string
	Options     string
	Description string
	Required    int
}

var Setting = &setting{}

func (c setting) GetPublicSettings() (map[string]interface{}, error) {
	var returnData = make(map[string]interface{})
	var data []setting
	result := database.Connect(&setting{}).Where("private != ?", TypeIsPrivate).Find(&data)
	if result.Error != nil {
		return returnData, errors.New("获取配置失败")
	}
	for _,value := range data{
		returnData[value.Key] = value.Value
	}
	return returnData, nil
}
