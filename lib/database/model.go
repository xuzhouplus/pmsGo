package database

import (
	"gorm.io/gorm"
)

type Model interface {
	DB() *gorm.DB
}

func Connect(model Model) *gorm.DB {
	return model.DB()
}

func Query(model interface{}) *gorm.DB {
	return DB.Model(model)
}
