package model

import (
	"errors"
	"log"
	"pmsGo/lib/database"
	"pmsGo/lib/security"
	"time"
)

const (
	StatusEnabled  = 1
	StatusDisabled = 2
)

type Admin struct {
	database.Model
	ID        uint      `gorm:"primarykey" json:"id"`
	Uuid      string    `gorm:"unique" json:"uuid"`
	Type      string    `json:"type"`
	Account   string    `gorm:"uniqueIndex" json:"account"`
	Avatar    string    `json:"avatar"`
	Password  string    `json:"_"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Salt      string    `json:"_"`
}

var AdminModel = &Admin{}

func (model Admin) ValidatePassword(inputPassword string) (bool, error) {
	password, err := security.RsaDecryptByPrivateKey(inputPassword)
	if err != nil {
		log.Printf("password:%e \n", err)
		return false, err
	}
	password = security.MD5(password, model.Salt)
	if password == model.Password {
		return true, nil
	}
	return false, errors.New("password wrong")
}

func (model *Admin) SetPassword(inputPassword string) {
	salt := security.Uuid(true)
	password := security.MD5(inputPassword, salt)
	model.Salt = salt
	model.Password = password
}

func (model Admin) Login(account string, password string) (Admin, error) {
	var login Admin
	if account == "" || password == "" {
		return login, errors.New("登录失败")
	}
	result := database.Query(&Admin{}).Where("account = ?", account).Take(&login)
	if result.Error != nil {
		return login, errors.New("登录失败")
	}
	validate, err := login.ValidatePassword(password)
	if validate {
		return login, nil
	}
	return login, err
}

