package gateway

import (
	"fmt"
	"github.com/idoubi/goz"
	"log"
	"net/url"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
)

const FacebookGatewayType = "facebook"

const FacebookScopeType = "public_profile"
const FacebookGrantType = "authorization_code"
const (
	FacebookAuthorizeUrl   = "https://www.facebook.com/v10.0/dialog/oauth"
	FacebookAccessTokenUrl = "https://graph.facebook.com/v10.0/oauth/access_token"
	FacebookAccessUserUrl  = "https://graph.facebook.com/me"
)

type FacebookAccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
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

func (gateway Facebook) AuthorizeUrl(scope string, redirect string, state string) (string, string, error) {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", gateway.FacebookAppId)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	queryString := query.Encode()
	return FacebookAuthorizeUrl + "?" + queryString, state, nil
}

func (gateway Facebook) AccessToken(callbackData map[string]string, redirect string) (string, error) {
	client := goz.NewClient()
	response, err := client.Get(FacebookAccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"client_id":     gateway.FacebookAppId,
			"client_secret": gateway.FacebookAppSecret,
			"code":          callbackData["code"],
			"redirect_uri":  redirect,
		},
	})
	if err != nil {
		return "", err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return "", err
	}
	return body.Get("access_token").String(), nil
}

func (gateway Facebook) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(FacebookAccessUserUrl, goz.Options{
		Query: map[string]string{
			"access_token": accessToken,
			"fields":       "id,name,gender,picture.width(400)",
		},
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	log.Println(body)
	sex := "2"
	if body.Get("gender").String() == "main" {
		sex = "1"
	}
	return map[string]string{
		"avatar":   body.Get("picture.data.url").String(),
		"channel":  "0",
		"nickname": body.Get("name").String(),
		"gender":   sex,
		"open_id":  body.Get("id").String(),
		"union_id": body.Get("id").String(),
	}, nil
}
