package model

import (
	"pmsGo/lib/config"
)

const (
	SettingCarouselLimit    = "carousel_limit" //站点名
	SettingKeyTitle         = "title"          //站点名
	SettingKeyIcp           = "icp"            //备案号
	SettingKeyVersion       = "version"        //版本
	SettingKeyMaintain      = "maintain"       //维护状态
	SettingKeyIcon          = "icon"
	SettingKeyLogo          = "logo"
	SettingKeyEncryptSecret = "encrypt_secret" //加密密钥
	SettingKeyLoginDuration = "login_duration" //登录有效时长
)

const (
	MaintainTrue  = "true"
	MaintainFalse = "false"
)

type SiteSetting struct {
	Setting
}

var SiteSettingModel = &SiteSetting{}

func (model SiteSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model SiteSetting) Keys() []string {
	return []string{SettingCarouselLimit, SettingKeyTitle, SettingKeyIcp, SettingKeyVersion, SettingKeyMaintain, SettingKeyIcon, SettingKeyLogo, SettingKeyEncryptSecret, SettingKeyLoginDuration}
}
