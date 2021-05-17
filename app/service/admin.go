package service

import (
	"errors"
	"pmsGo/app/model"
	"pmsGo/lib/database"
)

type Admin struct {
}

var AdminService = &Admin{}

func (service Admin) findOneByAccount(account string, status int) (*model.Admin, error) {
	admin := &model.Admin{}
	connect := database.Query(&model.Admin{})
	connect.Where("account = ?", account)
	if status != 0 {
		connect.Where("status = ?", status)
	}
	connect.Limit(1)
	result := connect.Find(&admin)
	if result.Error != nil {
		return nil, result.Error
	}
	return admin, nil
}

func (service Admin) Login(account string, password string) (*model.Admin, error) {
	if account == "" || password == "" {
		return nil, errors.New("登录失败")
	}
	admin, err := service.findOneByAccount(account, model.AdminStatusEnabled)
	if err != nil {
		return nil, errors.New("登录失败")
	}
	validate, err := admin.ValidatePassword(password)
	if validate {
		return nil, nil
	}
	return admin, err
}
