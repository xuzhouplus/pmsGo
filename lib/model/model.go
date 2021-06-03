package model

import (
	"gorm.io/gorm"
	"pmsGo/lib/database"
)

type Model interface {
	DB() *gorm.DB
}

func Connect(model Model) *gorm.DB {
	return model.DB()
}

func Query(model interface{}) *gorm.DB {
	return database.DB.Model(model)
}
