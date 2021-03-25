package model

import "pmsGo/lib/config"

const (
	SettingKeyBaiduApiKey    = "baidu_api_key"
	SettingKeyBaiduSecretKey = "baidu_secret_key"
)

type BaiduSetting struct {
	Setting
}

var BaiduSettingModel = &BaiduSetting{}

func (model BaiduSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model BaiduSetting) Keys() []string {
	return []string{SettingKeyBaiduApiKey, SettingKeyBaiduSecretKey}
}
