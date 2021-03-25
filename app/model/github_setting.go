package model

import "pmsGo/lib/config"

const (
	SettingKeyGithubApplicationName = "github_application_name"
	SettingKeyGithubAppId           = "github_app_id"
	SettingKeyGithubAppSecret       = "github_app_secret"
)

type GithubSetting struct {
	Setting
}

var GithubSettingModel = &GithubSetting{}

func (model GithubSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model GithubSetting) Keys() []string {
	return []string{SettingKeyGithubAppId, SettingKeyGithubAppSecret, SettingKeyGithubApplicationName}
}
