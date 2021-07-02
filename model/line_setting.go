package model

import "pmsGo/lib/config"

const (
	SettingKeyLineAppId     = "line_app_id"
	SettingKeyLineAppSecret = "line_app_secret"
)

type LineSetting struct {
	Setting
}

var LineSettingModel = &LineSetting{}

func (model LineSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model LineSetting) Keys() []string {
	return []string{SettingKeyLineAppId, SettingKeyLineAppSecret}
}
