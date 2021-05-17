package model

import (
	"pmsGo/lib/config"
)

type CarouselSetting struct {
	Setting
}

var CarouselSettingModel = &CarouselSetting{}

const (
	SettingKeyCarouselType     = "carousel_type"     //轮播类型
	SettingKeyCarouselLimit    = "carousel_limit"    //轮播数量限制
	SettingKeyCarouselInterval = "carousel_interval" //轮播间隔时间
)

func (model CarouselSetting) TableName() string {
	return config.Config.Database.Prefix + "settings"
}
func (model CarouselSetting) Keys() []string {
	return []string{SettingKeyCarouselType, SettingKeyCarouselLimit, SettingKeyCarouselInterval}
}
