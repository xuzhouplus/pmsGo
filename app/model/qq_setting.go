package model

import "pmsGo/lib/config"

const (
	SettingKeyQqAppId     = "qq_app_id"
	SettingKeyQqAppSecret = "qq_app_secret"
)

type QqSetting struct {
	Setting
}

var QqSettingModel = &QqSetting{}

func (model QqSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model QqSetting) Keys() []string {
	return []string{SettingKeyQqAppId, SettingKeyQqAppSecret}
}
