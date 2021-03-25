package model

import "pmsGo/lib/config"

const (
	SettingKeyFacebookAppId     = "facebook_app_id"
	SettingKeyFacebookAppSecret = "facebook_app_secret"
)

type FacebookSetting struct {
	Setting
}

var FacebookSettingModel = &FacebookSetting{}

func (model FacebookSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model FacebookSetting) Keys() []string {
	return []string{SettingKeyFacebookAppId, SettingKeyFacebookAppSecret}
}
