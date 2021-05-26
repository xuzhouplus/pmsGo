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

const WechatGatewayType = "wechat"
const WechatScopeType = "snsapi_login"
const WechatGrantType = "authorization_code"
const (
	WechatAuthorizeUrl   = "https://open.weixin.qq.com/connect/qrconnect"
	WechatAccessTokenUrl = "https://api.weixin.qq.com/sns/oauth2/access_token"
	WechatAccessUserUrl  = "https://api.weixin.qq.com/sns/userinfo"
)

type WechatAccessTokenRequest struct {
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
	Code      string `json:"code"`
	GrantType string `json:"grant_type"`
}

type Wechat struct {
	WechatOpenId    string
	WechatAppId     string
	WechatAppSecret string
}

func NewWechat() (*Wechat, error) {
	gateway := &Wechat{}
	appId := service.SettingService.GetSetting(model.SettingKeyWechatAppId)
	if appId == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyWechatAppId)
	}
	gateway.WechatAppId = appId
	appSecret := service.SettingService.GetSetting(model.SettingKeyWechatAppSecret)
	if appSecret == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyWechatAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appSecret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gateway.WechatAppSecret = string(decrypt)
	return gateway, nil
}

func (gateway Wechat) Scope() string {
	return WechatScopeType
}

func (gateway Wechat) GrantType() string {
	return WechatGrantType
}

func (gateway Wechat) AuthorizeUrl(scope string, redirect string, state string) (string, string, error) {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("appid", gateway.WechatAppId)
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	query.Add("response_type", "code")
	queryString := query.Encode()
	return WechatAuthorizeUrl + "?" + queryString + "#wechat_redirect", state, nil
}

func (gateway *Wechat) AccessToken(callbackData map[string]string, redirect string) (string, error) {
	requestData := &WechatAccessTokenRequest{
		Appid:     gateway.WechatAppId,
		Secret:    gateway.WechatAppSecret,
		GrantType: gateway.GrantType(),
		Code:      callbackData["code"],
	}
	client := goz.NewClient()
	response, err := client.Post(WechatAccessTokenUrl, goz.Options{
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
	gateway.WechatOpenId = body.Get("openid").String()
	return body.Get("access_token").String(), nil
}

func (gateway Wechat) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(WechatAccessUserUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"access_token": accessToken,
			"openid":       gateway.WechatOpenId,
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
	if body.Get("gender").String() == "1" {
		sex = "1"
	}
	return map[string]string{
		"avatar":   body.Get("headimgurl").String(),
		"channel":  "0",
		"nickname": body.Get("nickname").String(),
		"gender":   sex,
		"open_id":  gateway.WechatOpenId,
		"union_id": body.Get("unionid").String(),
	}, nil
}
