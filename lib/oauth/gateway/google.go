package gateway

import (
	"fmt"
	"github.com/idoubi/goz"
	"net/url"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
	model2 "pmsGo/model"
	service2 "pmsGo/service"
)

const GoogleGatewayType = "google"
const GoogleScopeType = "https://www.googleapis.com/auth/userinfo.profile"
const GoogleGrantType = "authorization_code"
const (
	GoogleAuthorizeUrl   = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleAccessTokenUrl = "https://oauth2.googleapis.com/token"
	GoogleAccessUserUrl  = "https://accounts.google.com/o/oauth2/v2/userinfo"
)

type GoogleAccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	RedirectUri  string `json:"redirect_uri"`
}

type Google struct {
	GoogleAppId     string
	GoogleAppSecret string
}

func NewGoogle() (*Google, error) {
	gateway := &Google{}
	appId := service2.SettingService.GetSetting(model2.SettingKeyGoogleAppId)
	if appId == "" {
		return nil, fmt.Errorf("缺少配置：%v", model2.SettingKeyGoogleAppId)
	}
	gateway.GoogleAppId = appId
	appSecret := service2.SettingService.GetSetting(model2.SettingKeyGoogleAppSecret)
	if appSecret == "" {
		return nil, fmt.Errorf("缺少配置：%v", model2.SettingKeyGoogleAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appSecret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gateway.GoogleAppSecret = string(decrypt)
	return gateway, nil
}

func (gateway Google) Scope() string {
	return GoogleScopeType
}

func (gateway Google) GrantType() string {
	return GoogleGrantType
}

func (gateway Google) AuthorizeUrl(scope string, redirect string, state string) (string, string, error) {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", gateway.GoogleAppId)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	queryString := query.Encode()
	return GoogleAuthorizeUrl + "?" + queryString, state, nil
}

func (gateway Google) AccessToken(callbackData map[string]string, redirect string) (map[string]string, error) {
	requestData := &GoogleAccessTokenRequest{
		ClientId:     gateway.GoogleAppId,
		ClientSecret: gateway.GoogleAppSecret,
		Code:         callbackData["code"],
		GrantType:    gateway.GrantType(),
		RedirectUri:  redirect,
	}
	client := goz.NewClient()
	response, err := client.Post(GoogleAccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		JSON: requestData,
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"accessToken":  body.Get("access_token").String(),
		"refreshToken": body.Get("access_token").String(),
	}, nil
}

func (gateway Google) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(GoogleAccessUserUrl, goz.Options{
		Headers: map[string]interface{}{
			"Authorization": "Bearer " + accessToken,
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
	if body.Get("gender").String() == "male" {
		sex = "1"
	}
	return map[string]string{
		"avatar":   body.Get("picture").String(),
		"channel":  "0",
		"nickname": body.Get("name").String(),
		"gender":   sex,
		"open_id":  body.Get("id").String(),
		"union_id": body.Get("id").String(),
	}, nil
}
