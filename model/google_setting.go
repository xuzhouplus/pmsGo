package model

import "pmsGo/lib/config"

const (
	SettingKeyGoogleAppId     = "google_app_id"
	SettingKeyGoogleAppSecret = "google_app_secret"
)

type GoogleSetting struct {
	Setting
}

var GoogleSettingModel = &GoogleSetting{}

func (model GoogleSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model GoogleSetting) Keys() []string {
	return []string{SettingKeyGoogleAppId, SettingKeyGoogleAppSecret}
}
