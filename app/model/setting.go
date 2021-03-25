package model

import (
	"errors"
	"pmsGo/lib/config"
	"pmsGo/lib/database"
	"pmsGo/lib/helper"
	"reflect"
)

const TypeIsPrivate = 1

type Setting struct {
	database.Model
	ID          int    `gorm:"private_key" json:"id"`
	Key         string `gorm:"uniqueIndex" json:"key"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Private     int    `json:"private"`
	Value       string `json:"value"`
	Options     string `json:"options"`
	Description string `json:"description"`
	Required    int    `json:"required"`
}

var SettingModel = &Setting{}

var Settings map[string]interface{}

func GetSettings(key string) interface{} {
	if len(Settings) == 0 {
		Settings = SettingModel.GetSettings([]string{})
	}
	if key != "" {
		return Settings[key]
	}
	return Settings
}

func (model Setting) GetPublicSettings() (map[string]interface{}, error) {
	var returnData = make(map[string]interface{})
	var data []Setting
	result := database.Query(&Setting{}).Where("private != ?", TypeIsPrivate).Find(&data)
	if result.Error != nil {
		return returnData, errors.New("获取配置失败")
	}
	for _, value := range data {
		returnData[value.Key] = value.Value
	}
	returnData["connects"] = config.Config.Web["connects"]
	return returnData, nil
}

func (model Setting) GetSetting(key string) string {
	var record Setting
	result := database.Query(&Setting{}).Where("`key` = ?", key).Limit(1).Take(&record)
	if result.Error != nil {
		return ""
	}
	return record.Value
}

func (model Setting) GetSettings(keys []string) map[string]interface{} {
	var data []Setting
	var ret = make(map[string]interface{})
	query := database.Query(&Setting{}).Select("key", "value")
	if len(keys) > 0 {
		query.Where("`key` IN (?)", keys).Limit(len(keys))
	}
	result := query.Find(&data)
	if result.Error != nil {
		return ret
	}
	if len(data) == 0 {
		return ret
	}
	for _, record := range data {
		ret[record.Key] = record.Value
	}
	return ret
}
func (model Setting) Find(keys []string, indexBy string) map[string]Setting {
	var data []Setting
	result := database.Query(&Setting{}).Where("`key` IN ?", keys).Limit(len(keys)).Find(&data)
	if result.Error != nil {
		return nil
	}
	if result.RowsAffected == 0 {
		return nil
	}
	var list = make(map[string]Setting)
	for _, setting := range data {
		e := reflect.ValueOf(&setting).Elem()
		field := e.FieldByName(helper.FirstToUpper(indexBy))
		list[field.Interface().(string)] = setting
	}
	return list
}
func (model Setting) Save(keyPairs map[string]interface{}) {
	for key, value := range keyPairs {
		database.Query(&Setting{}).Where("`key` = ?", key).Update("value", value)
	}
}
