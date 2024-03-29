package model

import "pmsGo/lib/config"

const (
	SettingKeyBaiduAppName         = "baidu_app_name"
	SettingKeyBaiduApiKey          = "baidu_api_key"
	SettingKeyBaiduSecretKey       = "baidu_secret_key"
	SettingKeyBaiduPanAvailability = "baidu_pan_availability"
)

const (
	BaiduPanDisabled = "disabled"
	BaiduPanEnabled  = "enabled"
)

type BaiduSetting struct {
	Setting
}

var BaiduSettingModel = &BaiduSetting{}

func (model BaiduSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model BaiduSetting) Keys() []string {
	return []string{SettingKeyBaiduAppName, SettingKeyBaiduApiKey, SettingKeyBaiduSecretKey, SettingKeyBaiduPanAvailability}
}
