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

const WeiboGatewayType = "weibo"

const WeiboScopeType = "all"
const WeiboGrantType = "authorization_code"
const WeiboAuthorizeDisplay = "default"
const (
	WeiboAuthorizeUrl   = "https://api.weibo.com/oauth2/authorize"
	WeiboAccessTokenUrl = "https://api.weibo.com/oauth2/access_token"
	WeiboAccessUserUrl  = "https://api.weibo.com/2/users/show.json"
)

type WeiboAccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RedirectUri  string `json:"redirect_uri"`
	Code         string `json:"code"`
}
type WeiboAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	RemindIn    string `json:"remind_in"`
	Uid         string `json:"uid"`
}
type WeiboAccessUserRequest struct {
	AccessToken string `json:"access_token"`
	Uuid        string `json:"uuid"`
}
type Weibo struct {
	WeiboUuid      string
	WeiboAppKey    string
	WeiboAppSecret string
}

func NewWeibo() (*Weibo, error) {
	gateway := &Weibo{}
	appId := service.SettingService.GetSetting(model.SettingKeyWeiboAppId)
	if appId == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyWeiboAppId)
	}
	gateway.WeiboAppKey = appId
	appSecret := service.SettingService.GetSetting(model.SettingKeyWeiboAppSecret)
	if appSecret == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyWeiboAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appSecret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gateway.WeiboAppSecret = string(decrypt)
	return gateway, nil
}

func (gateway Weibo) Scope() string {
	return WeiboScopeType
}

func (gateway Weibo) GrantType() string {
	return WeiboGrantType
}

func (gateway Weibo) AuthorizeUrl(scope string, redirect string, state string) string {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", gateway.WeiboAppKey)
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	query.Add("display", WeiboAuthorizeDisplay)
	query.Add("forcelogin", "true")
	queryString := query.Encode()
	return WeiboAuthorizeUrl + "?" + queryString
}

func (gateway *Weibo) AccessToken(code string, redirect string, state string) (string, error) {
	requestData := &WeiboAccessTokenRequest{
		ClientId:     gateway.WeiboAppKey,
		ClientSecret: gateway.WeiboAppSecret,
		GrantType:    gateway.GrantType(),
		Code:         code,
		RedirectUri:  redirect,
	}
	client := goz.NewClient()
	response, err := client.Post(WeiboAccessTokenUrl, goz.Options{
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
	gateway.WeiboUuid = body.Get("uuid").String()
	return body.Get("access_token").String(), nil
}

func (gateway Weibo) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(WeiboAccessUserUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"access_token": accessToken,
			"uuid":         gateway.WeiboUuid,
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
	if body.Get("gender").String() == "m" {
		sex = "1"
	}
	return map[string]string{
		"avatar":   body.Get("avatar_hd").String(),
		"channel":  "0",
		"nickname": body.Get("screen_name").String(),
		"gender":   sex,
		"open_id":  gateway.WeiboUuid,
		"union_id": gateway.WeiboUuid,
	}, nil
}
