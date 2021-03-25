package database

import (
	"gorm.io/gorm"
)

type Model struct {
}

func Query(model interface{}) *gorm.DB {
	return DB.Model(model)
}
