package model

import (
	"pmsGo/lib/config"
)

const (
	SettingKeyTitle         = "title"    //站点名
	SettingKeyIcp           = "icp"      //备案号
	SettingKeyVersion       = "version"  //版本
	SettingKeyMaintain      = "maintain" //维护状态
	SettingKeyIcon          = "icon"
	SettingKeyLogo          = "logo"
	SettingKeyEncryptSecret = "encrypt_secret" //加密密钥
	SettingKeyLoginDuration = "login_duration" //登录有效时长
)

type SiteSetting struct {
	Setting
}

var SiteSettingModel = &SiteSetting{}

func (model SiteSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model SiteSetting) Keys() []string {
	return []string{SettingKeyTitle, SettingKeyIcp, SettingKeyVersion, SettingKeyMaintain, SettingKeyIcon, SettingKeyLogo, SettingKeyEncryptSecret, SettingKeyLoginDuration}
}

