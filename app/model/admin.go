package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"pmsGo/lib/database"
	"pmsGo/lib/security/md5"
	"pmsGo/lib/security/random"
	"pmsGo/lib/security/rsa"
	"time"
)

const (
	AdminStatusEnabled  = 1
	AdminStatusDisabled = 2
)

type Admin struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Uuid      string    `gorm:"unique" json:"uuid"`
	Type      string    `json:"type"`
	Account   string    `gorm:"uniqueIndex" json:"account"`
	Avatar    string    `json:"avatar"`
	Password  string    `json:"_"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Salt      string    `json:"_"`
}

func (model Admin) DB() *gorm.DB {
	return database.DB.Model(&model)
}

func (model Admin) ValidatePassword(inputPassword string) (bool, error) {
	password, err := rsa.DecryptByPrivateKey(inputPassword)
	if err != nil {
		log.Printf("password:%e \n", err)
		return false, err
	}

	password = md5.Md5(password, model.Salt)
	if password == model.Password {
		return true, nil
	}
	return false, errors.New("password wrong")
}

func (model *Admin) SetPassword(inputPassword string) {
	salt := random.Uuid(true)
	password := md5.Md5(inputPassword, salt)
	model.Salt = salt
	model.Password = password
}
