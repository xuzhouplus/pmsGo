package model

import "pmsGo/lib/config"

const (
	SettingKeyGiteeApplicationName = "gitee_application_name"
	SettingKeyGiteeAppId           = "gitee_app_id"
	SettingKeyGiteeAppSecret       = "gitee_app_secret"
)

type GiteeSetting struct {
	Setting
}

var GiteeSettingModel = &GiteeSetting{}

func (model GiteeSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model GiteeSetting) Keys() []string {
	return []string{SettingKeyGiteeAppId, SettingKeyGiteeAppSecret, SettingKeyGiteeApplicationName}
}
