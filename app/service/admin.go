package service

import (
	"errors"
	"fmt"
	"log"
	"pmsGo/app/model"
	"pmsGo/lib/oauth/user"
	"pmsGo/lib/security/rsa"
)

type Admin struct {
}

var AdminService = &Admin{}

func (service Admin) FindOneByAccount(account string, status int) (*model.Admin, error) {
	admin := &model.Admin{}
	connect := admin.DB()
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
	connect := admin.DB()
	connect.Where("uuid = ?", uuid)
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

func (service Admin) FindOneById(id int, status int) (*model.Admin, error) {
	admin := &model.Admin{}
	connect := admin.DB()
	connect.Where("id = ?", id)
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
	admin, err := service.FindOneByAccount(account, model.AdminStatusEnabled)
	if err != nil {
		return nil, errors.New("登录失败")
	}
	validate, err := admin.ValidatePassword(password)
	if !validate {
		return nil, err
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
	avatar := postData["avatar"].(string)
	if avatar != "" {
		admin.Avatar = avatar
	}
	password := postData["password"].(string)
	if password != "" {
		password, err := rsa.DecryptByPrivateKey(password)
		if err != nil {
			log.Printf("password:%e \n", err)
			return nil, err
		}
		admin.SetPassword(password)
	}
	admin.Status = int(postData["status"].(float64))
	result := admin.DB().Save(&admin)
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
	connect := admin.DB()
	result := connect.Where("admin_id = ?", admin.ID).Find(&connects)
	if result.Error != nil {
		return nil, result.Error
	}
	return connects, nil
}

func (service Admin) Bind(adminId int, authAccount *user.User) (*model.Connect, error) {
	connect := &model.Connect{}
	query := connect.DB()
	result := query.Where("union_id = ?", authAccount.OpenId).Limit(1).Find(&connect)
	if result != nil {
		admin, err := service.FindOneById(connect.AdminId, model.AdminStatusEnabled)
		if err != nil {
			return nil, err
		}
		if admin == nil {
			return nil, fmt.Errorf("重复绑定账号：%v", admin.Account)
		}
	}
	connect.Account = authAccount.Nickname
	connect.UnionId = authAccount.OpenId
	connect.AdminId = adminId
	connect.Avatar = authAccount.Avatar
	connect.Type = authAccount.Type
	connect.Status = model.ConnectStatusEnable
	result = query.Create(&connect)
	if result.Error != nil {
		return nil, result.Error
	}
	return connect, nil
}

func (service Admin) Auth(authAccount *user.User) (*model.Admin, error) {
	connect := &model.Connect{}
	query := connect.DB()
	result := query.Where("union_id = ?", authAccount.OpenId).Where("status = ?", model.ConnectStatusEnable).Limit(1).Find(&connect)
	if result.Error != nil {
		return nil, result.Error
	}
	if connect == nil {
		return nil, fmt.Errorf("没有绑定账号：%v", authAccount.Nickname)
	}
	admin, err := service.FindOneById(connect.AdminId, model.AdminStatusEnabled)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
