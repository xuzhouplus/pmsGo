package database

import (
	"gorm.io/gorm"
	"pmsGo/lib/config"
)

type Model struct {
}

func Query(model interface{}) *gorm.DB {
	if config.Config.Site.Debug {
		return DB.Debug().Model(model)
	}
	return DB.Model(model)
}
