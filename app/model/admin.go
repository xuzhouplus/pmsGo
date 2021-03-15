package model

import (
	"errors"
	"fmt"
	"pmsGo/lib/database"
	"pmsGo/lib/security"
	"time"
)

const (
	StatusEnabled  = 1
	StatusDisabled = 2
)

type admin struct {
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

var Admin = &admin{}

func (model admin) ValidatePassword(inputPassword string) (bool, error) {
	password, err := security.RsaDecryptByPrivateKey(inputPassword)
	if err != nil {
		fmt.Printf("password:%e", err)
		return false, err
	}
	password = security.MD5(inputPassword, model.Salt)
	return password == model.Password, nil
}

func (model *admin) SetPassword(inputPassword string) {
	salt := security.Uuid(true)
	password := security.MD5(inputPassword, salt)
	model.Salt = salt
	model.Password = password
}

func (model admin) Login(account string, password string) (interface{}, error) {
	var login admin
	result := database.Query(&admin{}).Where("account = ? AND status = ?", account, StatusEnabled).First(&login)
	if result.Error != nil {
		return login, errors.New("登录失败")
	}
	validate, err := login.ValidatePassword(password)
	if validate {
		return login, nil
	}
	return nil, err
}

func (model admin) Logout()  {

}
