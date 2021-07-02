package model

import "pmsGo/lib/config"

const (
	SettingKeyTwitterAppId     = "twitter_app_id"
	SettingKeyTwitterAppSecret = "twitter_app_secret"
)

type TwitterSetting struct {
	Setting
}

var TwitterSettingModel = &TwitterSetting{}

func (model TwitterSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model TwitterSetting) Keys() []string {
	return []string{SettingKeyTwitterAppId, SettingKeyTwitterAppSecret}
}
