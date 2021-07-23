package service

import (
	"errors"
	"pmsGo/lib/config"
	"pmsGo/lib/helper"
	"pmsGo/model"
	"reflect"
	"strconv"
)

type Setting struct {
	Settings map[string]string
}

var SettingService = NewSettingService()

func NewSettingService() *Setting {
	service := &Setting{}
	settingModel := &model.Setting{}
	var data []model.Setting
	result := settingModel.DB().Find(&data)
	if result.Error != nil {
		panic(result.Error)
	}
	settings := make(map[string]string, 0)
	for _, value := range data {
		settings[value.Key] = value.Value
	}
	service.Settings = settings
	return service
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

func (service Setting) GetLoginSettings() []string {
	connects := config.Config.Web.Connects
	if len(connects) == 0 {
		return nil
	}
	connectSettings := make(map[string][]string)
	connectSettingKeys := make([]string, 0)
	for _, connect := range connects {
		settingKeys := make([]string, 0)
		switch connect {
		case model.ConnectTypeAlipay:
			settingKeys = model.AlipaySettingModel.Keys()
		case model.ConnectTypeBaidu:
			settingKeys = model.BaiduSettingModel.Keys()
		case model.ConnectTypeFacebook:
			settingKeys = model.FacebookSettingModel.Keys()
		case model.ConnectTypeGitee:
			settingKeys = model.GiteeSettingModel.Keys()
		case model.ConnectTypeGithub:
			settingKeys = model.GithubSettingModel.Keys()
		case model.ConnectTypeGoogle:
			settingKeys = model.GoogleSettingModel.Keys()
		case model.ConnectTypeLine:
			settingKeys = model.LineSettingModel.Keys()
		case model.ConnectTypeTwitter:
			settingKeys = model.TwitterSettingModel.Keys()
		case model.ConnectTypeQq:
			settingKeys = model.QqSettingModel.Keys()
		case model.ConnectTypeWechat:
			settingKeys = model.WechatSettingModel.Keys()
		case model.ConnectTypeWeibo:
			settingKeys = model.WeiboSettingModel.Keys()
		}
		connectSettingKeys = append(connectSettingKeys, settingKeys...)
		connectSettings[connect] = settingKeys
	}
	connectSettingVals := service.GetSettings(connectSettingKeys)
	available := make([]string, 0)
	for connectType, settings := range connectSettings {
		empty := false
		for _, setting := range settings {
			if connectSettingVals[setting] == "" {
				empty = true
				break
			}
		}
		if !empty {
			available = append(available, connectType)
		}
	}
	return available
}

func (service Setting) GetPublicSettings() (map[string]interface{}, error) {
	var returnData = make(map[string]interface{})
	var data []model.Setting
	settingModel := &model.Setting{}
	result := settingModel.DB().Where("private != ?", model.SettingTypeIsPrivate).Find(&data)
	if result.Error != nil {
		return nil, errors.New("获取配置失败")
	}
	for _, value := range data {
		returnData[value.Key] = value.Value
	}
	return returnData, nil
}

func (service Setting) GetSetting(key string) string {
	if service.Settings != nil {
		return service.Settings[key]
	}
	settingModel := &model.Setting{}
	result := settingModel.DB().Where("`key` = ?", key).Limit(1).Take(&settingModel)
	if result.Error != nil {
		return ""
	}
	return settingModel.Value
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
	settingModel := &model.Setting{}
	query := settingModel.DB().Select("key", "value")
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
	settingModel := &model.Setting{}
	result := settingModel.DB().Where("`key` IN ?", keys).Limit(len(keys)).Find(&data)
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
	settingModel := &model.Setting{}
	connect := settingModel.DB().Begin()
	for key, value := range keyPairs {
		if value == "" && service.IsUnsafeSetting(key) {
			continue
		}
		result := settingModel.DB().Where("`key` = ?", key).Update("value", value)
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

func (service Setting) BaiduPanAvailability() bool {
	baiduSettings := service.GetSettings(model.BaiduSettingModel.Keys())
	if baiduSettings[model.SettingKeyBaiduPanAvailability] == model.BaiduPanDisabled {
		return false
	}
	if baiduSettings[model.SettingKeyBaiduApiKey] == "" || baiduSettings[model.SettingKeyBaiduSecretKey] == "" {
		return false
	}
	return true
}
