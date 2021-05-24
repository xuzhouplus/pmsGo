package gateway

import (
	"fmt"
	"github.com/idoubi/goz"
	"net/url"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
)

const QqGatewayType = "qq"
const QqScopeType = "get_user_info"
const QqGrantType = ""
const (
	QqAuthorizeUrl   = "https://graph.qq.com/oauth2.0/authorize"
	QqAccessTokenUrl = "https://graph.qq.com/oauth2.0/token"
	QqAccessMeUrl    = "https://graph.qq.com/oauth2.0/me"
	QqAccessUserUrl  = "https://graph.qq.com/user/get_user_info"
)

type QqAccessTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
	Fmt          string `json:"fmt"`
}

type Qq struct {
	QqOpenId    string
	QqAppId     string
	QqAppSecret string
}

func NewQq() (*Qq, error) {
	gateway := &Qq{}
	appId := service.SettingService.GetSetting(model.SettingKeyQqAppId)
	if appId == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyQqAppId)
	}
	gateway.QqAppId = appId
	appSecret := service.SettingService.GetSetting(model.SettingKeyQqAppSecret)
	if appSecret == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyWeiboAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appSecret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gateway.QqAppSecret = string(decrypt)
	return gateway, nil
}

func (gateway Qq) Scope() string {
	return QqScopeType
}

func (gateway Qq) GrantType() string {
	return QqGrantType
}

func (gateway Qq) AuthorizeUrl(scope string, redirect string, state string) string {
	if scope == "" {
		scope = gateway.Scope()
	}
	url := url.URL{}
	query := url.Query()
	query.Add("client_id", gateway.QqAppId)
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	query.Add("response_type", "code")
	queryString := query.Encode()
	return QqAuthorizeUrl + "?" + queryString
}

func (gateway *Qq) AccessToken(code string, redirect string, state string) (string, error) {
	requestData := &QqAccessTokenRequest{
		GrantType:    gateway.GrantType(),
		ClientId:     gateway.QqAppId,
		ClientSecret: gateway.QqAppSecret,
		Code:         code,
		RedirectUri:  redirect,
		Fmt:          "json",
	}
	client := goz.NewClient()
	response, err := client.Post(QqAccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		JSON: requestData,
	})
	if err != nil {
		return "", err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return "", err
	}
	gateway.QqOpenId, err = gateway.Me(body.Get("access_token").String())
	if err != nil {
		return "", err
	}
	return body.Get("access_token").String(), nil
}

func (gateway Qq) Me(accessToken string) (string, error) {
	client := goz.NewClient()
	response, err := client.Get(QqAccessMeUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"access_token": accessToken,
		},
	})
	if err != nil {
		return "", err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return "", err
	}
	return body.Get("openid").String(), nil
}

func (gateway Qq) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(QqAccessUserUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"access_token":       accessToken,
			"oauth_consumer_key": gateway.QqAppId,
			"openid":             gateway.QqOpenId,
		},
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	sex := "2"
	if body.Get("gender").String() == "男" {
		sex = "1"
	}
	avatar := body.Get("figureurl_qq_2").String()
	if avatar == "" {
		avatar = body.Get("figureurl_qq_1").String()
	}
	return map[string]string{
		"avatar":   avatar,
		"channel":  "0",
		"nickname": body.Get("nickname").String(),
		"gender":   sex,
		"open_id":  gateway.QqOpenId,
		"union_id": gateway.QqOpenId,
	}, nil
}
