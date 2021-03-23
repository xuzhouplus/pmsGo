package model

import (
	"errors"
	"pmsGo/lib/database"
	"time"
)

type admin struct {
	database.Model
	ID        uint   `gorm:"primarykey"`
	Uuid      string `gorm:"unique"`
	Type      string
	Account   string `gorm:uniqueIndex`
	Avatar    string
	Password  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Salt      string
}

var Admin = &admin{}

func (c admin) Login(account string, password string) (admin, error) {
	var login admin
	result := database.Query(&admin{}).Where("account = ? AND password = ?", account, password).First(&login)
	if result.Error != nil {
		return login, errors.New("登录失败")
	}
	return login, nil
}
