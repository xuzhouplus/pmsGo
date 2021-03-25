package model

import "pmsGo/lib/config"

const (
	SettingKeyWeiboAppId     = "weibo_app_id"
	SettingKeyWeiboAppSecret = "weibo_app_secret"
)

type WeiboSetting struct {
	Setting
}

var WeiboSettingModel = &WeiboSetting{}

func (model WeiboSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model WeiboSetting) Keys() []string {
	return []string{SettingKeyWeiboAppId, SettingKeyWeiboAppSecret}
}
