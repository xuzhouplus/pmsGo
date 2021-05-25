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

const TwitterGatewayType = "twitter"

const TwitterScopeType = ""
const TwitterGrantType = ""
const (
	TwitterRequestTokenUrl = "https://api.twitter.com/oauth/request_token"
	TwitterAuthorizeUrl    = "https://api.twitter.com/oauth/authenticate"
	TwitterAccessTokenUrl  = "https://api.twitter.com/oauth/access_token"
	TwitterAccessUserUrl   = "https://api.twitter.com/1.1/account/verify_credentials.json"
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

func (gateway Twitter) RequestToken(redirect string) (string, error) {
	client := goz.NewClient()
	response, err := client.Post(BaiduAccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		FormParams: map[string]interface{}{
			"oauth_callback": redirect,
		},
	})
	if err != nil {
		return "", err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return "", err
	}
	return body.Get("oauth_token").String(), nil
}

func (gateway Twitter) AuthorizeUrl(scope string, redirect string, state string) (string, error) {
	token, err := gateway.RequestToken(redirect)
	if err != nil {
		return "", err
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("oauth_token", token)
	queryString := query.Encode()
	return TwitterAuthorizeUrl + "?" + queryString, nil
}

func (gateway Twitter) AccessToken(code string, redirect string, state string) (string, error) {

	requestData := &BaiduAccessTokenRequest{}
	client := goz.NewClient()
	response, err := client.Post(BaiduAccessTokenUrl, goz.Options{
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
	return body.Get("access_token").String(), nil
}

func (gateway Twitter) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(BaiduAccessUserUrl, goz.Options{
		Query: map[string]string{
			"access_token": accessToken,
			"get_unionid":  "1",
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
	return map[string]string{
		"avatar":   body.Get("portrait").String(),
		"channel":  "0",
		"nickname": body.Get("username").String(),
		"gender":   body.Get("sex").String(),
		"open_id":  body.Get("openid").String(),
		"union_id": body.Get("unionid").String(),
	}, nil
}
