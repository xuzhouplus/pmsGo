package model

import "pmsGo/lib/config"

const (
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
	return []string{SettingKeyBaiduApiKey, SettingKeyBaiduSecretKey, SettingKeyBaiduPanAvailability}
}
