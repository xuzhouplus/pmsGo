package model

import "pmsGo/lib/config"

const (
	SettingKeyWechatAppId     = "wechat_app_id"
	SettingKeyWechatAppSecret = "wechat_app_secret"
)

type WechatSetting struct {
	Setting
}

var WechatSettingModel = &WechatSetting{}

func (model WechatSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model WechatSetting) Keys() []string {
	return []string{SettingKeyWechatAppId, SettingKeyWechatAppSecret}
}
