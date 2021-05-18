package service

import (
	"errors"
	"log"
	"pmsGo/app/model"
	"pmsGo/lib/database"
	"pmsGo/lib/security"
)

type Admin struct {
}

var AdminService = &Admin{}

func (service Admin) FindOneByAccount(account string, status int) (*model.Admin, error) {
	admin := &model.Admin{}
	connect := database.Query(&model.Admin{})
	connect.Where("account = ?", account)
	if status != 0 {
		connect.Where("status = ?", status)
	}
	connect.Limit(1)
	result := connect.Take(&admin)
	if result.Error != nil {
		return nil, result.Error
	}
	return admin, nil
}

func (service Admin) FindOneByUuid(uuid string, status int) (*model.Admin, error) {
	admin := &model.Admin{}
	connect := database.Query(&model.Admin{})
	connect.Where("uuid = ?", uuid)
	if status != 0 {
		connect.Where("status = ?", status)
	}
	connect.Limit(1)
	result := connect.Take(&admin)
	if result.Error != nil {
		return nil, result.Error
	}
	return admin, nil
}

func (service Admin) Login(account string, password string) (*model.Admin, error) {
	if account == "" || password == "" {
		return nil, errors.New("登录失败")
	}
	admin, err := service.FindOneByAccount(account, model.AdminStatusEnabled)
	if err != nil {
		return nil, errors.New("登录失败")
	}
	validate, err := admin.ValidatePassword(password)
	if validate {
		return nil, nil
	}
	return admin, err
}

func (service Admin) Update(postData map[string]interface{}) (*model.Admin, error) {
	uuidData := postData["uuid"]
	if uuidData == nil {
		return nil, errors.New("账号唯一标识为空")
	}
	uuid := uuidData.(string)
	admin, err := service.FindOneByUuid(uuid, 0)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("账号不存在")
	}
	admin.Avatar = postData["avatar"].(string)
	password, err := security.RsaDecryptByPrivateKey(postData["password"].(string))
	if err != nil {
		log.Printf("password:%e \n", err)
		return nil, err
	}
	admin.SetPassword(password)
	admin.Status = int(postData["status"].(float64))
	connect := database.Query(&model.Admin{})
	result := connect.Where("uuid = ?", admin.Uuid).Save(&admin)
	if result.Error != nil {
		return nil, result.Error
	}
	return admin, nil
}

func (service Admin) GetBoundConnects(account string) ([]model.Connect, error) {
	admin, err := service.FindOneByAccount(account, 0)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("账号不存在")
	}
	var connects []model.Connect
	connect := database.Query(&model.Connect{})
	result := connect.Where("admin_id = ?", admin.ID).Find(&connects)
	if result.Error != nil {
		return nil, result.Error
	}
	return connects, nil
}
