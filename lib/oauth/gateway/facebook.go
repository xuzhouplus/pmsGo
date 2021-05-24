package gateway

import (
	"fmt"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
)

const FacebookGatewayType = "facebook"

const FacebookScopeType = ""
const FacebookGrantType = ""
const (
	FacebookAuthorizeUrl    = ""
	FacebookAccessTokenUrl  = ""
	FacebookAccessUserUrl = ""
)

type FacebookAccessTokenRequest struct {
}

type Facebook struct {
	FacebookAppId     string
	FacebookAppSecret string
}

func NewFacebook() (*Facebook, error) {
	gateway := &Facebook{}
	appId := service.SettingService.GetSetting(model.SettingKeyFacebookAppId)
	if appId == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyFacebookAppId)
	}
	gateway.FacebookAppId = appId
	appSecret := service.SettingService.GetSetting(model.SettingKeyFacebookAppSecret)
	if appSecret == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyFacebookAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appSecret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gateway.FacebookAppSecret = string(decrypt)
	return gateway, nil
}

func (gateway Facebook) Scope() string {
	return FacebookScopeType
}

func (gateway Facebook) GrantType() string {
	return FacebookGrantType
}

func (gateway Facebook) AuthorizeUrl(scope string, redirect string, state string) string {
	panic("implement me")
}

func (gateway Facebook) AccessToken(code string, redirect string, state string) (string, error) {
	panic("implement me")
}

func (gateway Facebook) User(accessToken string) (map[string]string, error) {
	panic("implement me")
}
