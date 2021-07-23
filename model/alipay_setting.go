package model

import "pmsGo/lib/config"

const (
	SettingKeyAlipayAppId           = "alipay_app_id"          //应用appid
	SettingKeyAlipayAppPrimaryKey   = "alipay_app_primary_key" //应用私钥
	SettingKeyAlipayPublicKay = "alipay_public_key"      //支付宝公钥name = 
)

type AlipaySetting struct {
	Setting
}

var AlipaySettingModel = &AlipaySetting{}

func (model AlipaySetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}

func (model AlipaySetting) Keys()[]string  {
	return []string{SettingKeyAlipayAppId, SettingKeyAlipayAppPrimaryKey, SettingKeyAlipayPublicKay}
}
