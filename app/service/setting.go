package service

import (
	"errors"
	"pmsGo/app/model"
	"pmsGo/lib/config"
	"pmsGo/lib/database"
	"pmsGo/lib/helper"
	"reflect"
	"strconv"
)

type Setting struct {
	Settings map[string]string
}

var SettingService = &Setting{}
var Settings map[string]string

func (service *Setting) Load() {
	var data []model.Setting
	result := database.Query(&model.Setting{}).Find(&data)
	if result.Error != nil {
		panic(result.Error)
	}
	settings := make(map[string]string, 0)
	for _, value := range data {
		settings[value.Key] = value.Value
	}
	service.Settings = settings
}
func (service Setting) UnsafeSettings() []string {
	return []string{
		model.SettingKeyBaiduSecretKey,
		model.SettingKeyAlipayAppPrimaryKey,
		model.SettingKeyTwitterAppSecret,
		model.SettingKeyWechatAppSecret,
		model.SettingKeyWeiboAppSecret,
		model.SettingKeyGithubAppSecret,
		model.SettingKeyGoogleAppSecret,
		model.SettingKeyTwitterAppSecret,
		model.SettingKeyLineAppSecret,
		model.SettingKeyFacebookAppSecret,
		model.SettingKeyGiteeAppSecret,
	}
}

func (service Setting) IsUnsafeSetting(settingKey string) bool {
	unsafeSettings := service.UnsafeSettings()
	for _, setting := range unsafeSettings {
		if setting == settingKey {
			return true
		}
	}
	return false
}

func (service Setting) GetPublicSettings() (map[string]interface{}, error) {
	var returnData = make(map[string]interface{})
	var data []model.Setting
	result := database.Query(&model.Setting{}).Where("private != ?", model.SettingTypeIsPrivate).Find(&data)
	if result.Error != nil {
		return nil, errors.New("获取配置失败")
	}
	for _, value := range data {
		returnData[value.Key] = value.Value
	}
	returnData["connects"] = config.Config.Web.Connects
	return returnData, nil
}

func (service Setting) GetSetting(key string) string {
	if service.Settings != nil {
		return service.Settings[key]
	}
	var record model.Setting
	result := database.Query(&model.Setting{}).Where("`key` = ?", key).Limit(1).Take(&record)
	if result.Error != nil {
		return ""
	}
	return record.Value
}
func (service *Setting) SetSetting(key string, value string) {
	if service.Settings == nil {
		service.Settings = make(map[string]string)
	}
	service.Settings[key] = value
}
func (service Setting) GetSettings(keys []string) map[string]string {
	var ret = make(map[string]string)
	if service.Settings != nil {
		for _, key := range keys {
			ret[key] = service.GetSetting(key)
		}
		return ret
	}
	var data []model.Setting
	query := database.Query(&model.Setting{}).Select("key", "value")
	if len(keys) > 0 {
		query.Where("`key` IN (?)", keys).Limit(len(keys))
	}
	result := query.Find(&data)
	if result.Error != nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	for _, record := range data {
		ret[record.Key] = record.Value
	}
	return ret
}
func (service Setting) Find(keys []string, indexBy string) map[string]model.Setting {
	var data []model.Setting
	result := database.Query(&model.Setting{}).Where("`key` IN ?", keys).Limit(len(keys)).Find(&data)
	if result.Error != nil {
		return nil
	}
	if result.RowsAffected == 0 {
		return nil
	}
	var list = make(map[string]model.Setting)
	for _, setting := range data {
		e := reflect.ValueOf(&setting).Elem()
		field := e.FieldByName(helper.FirstToUpper(indexBy))
		list[field.Interface().(string)] = setting
	}
	return list
}

func (service *Setting) Save(keyPairs map[string]interface{}) error {
	connect := database.DB.Begin()
	for key, value := range keyPairs {
		if value == "" && service.IsUnsafeSetting(key) {
			continue
		}
		result := connect.Model(&model.Setting{}).Where("`key` = ?", key).Update("value", value)
		if result.Error != nil {
			connect.Rollback()
			return result.Error
		}
	}
	for key, value := range keyPairs {
		if value == "" && service.IsUnsafeSetting(key) {
			continue
		}
		var setVal string
		switch value.(type) {
		case string:
			setVal = value.(string)
		case int:
			setVal = strconv.Itoa(value.(int))
		case int64:
			setVal = strconv.Itoa(int(value.(int64)))
		case float64:
			setVal = strconv.Itoa(int(value.(float64)))
		case bool:
			if value.(bool) {
				setVal = "true"
			} else {
				setVal = "false"
			}
		}
		service.SetSetting(key, setVal)
	}
	connect.Commit()
	return nil
}
