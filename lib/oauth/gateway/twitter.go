package gateway

import (
	"fmt"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
)

const TwitterGatewayType = "twitter"

const TwitterScopeType = ""
const TwitterGrantType = ""
const (
	TwitterAuthorizeUrl    = ""
	TwitterAccessTokenUrl  = ""
	TwitterAccessUserUrl = ""
)

type TwitterAccessTokenRequest struct {
}

type Twitter struct {
	TwitterAppId     string
	TwitterAppSecret string
}

func NewTwitter() (*Twitter, error) {
	gateway := &Twitter{}
	appId := service.SettingService.GetSetting(model.SettingKeyTwitterAppId)
	if appId == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyTwitterAppId)
	}
	gateway.TwitterAppId = appId
	appSecret := service.SettingService.GetSetting(model.SettingKeyTwitterAppSecret)
	if appSecret == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyTwitterAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appSecret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gateway.TwitterAppSecret = string(decrypt)
	return gateway, nil
}

func (gateway Twitter) Scope() string {
	return TwitterScopeType
}

func (gateway Twitter) GrantType() string {
	return TwitterGrantType
}

func (gateway Twitter) AuthorizeUrl(scope string, redirect string, state string) string {
	panic("implement me")
}

func (gateway Twitter) AccessToken(code string, redirect string, state string) (string, error) {
	panic("implement me")
}

func (gateway Twitter) User(accessToken string) (map[string]string, error) {
	panic("implement me")
}
