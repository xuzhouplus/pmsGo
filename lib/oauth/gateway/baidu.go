package gateway

import (
	"fmt"
	"github.com/idoubi/goz"
	"log"
	"net/url"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
	"pmsGo/model"
	"pmsGo/service"
)

const BaiduGatewayType = "baidu"
const BaiduUserScope = "basic"
const BaiduGrantType = "authorization_code"
const BaiduAuthorizeDisplay = "page"

const (
	BaiduAuthorizeUrl   = "https://openapi.baidu.com/oauth/2.0/authorize"
	BaiduAccessTokenUrl = "https://openapi.baidu.com/oauth/2.0/token"
	BaiduAccessUserUrl  = "https://openapi.baidu.com/rest/2.0/passport/users/getInf"
)

type BaiduAccessTokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

type BaiduAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type Baidu struct {
	BaiduApiKey    string
	BaiduSecretKey string
}

func NewBaidu() (*Baidu, error) {
	baiduGateway := &Baidu{}
	baiduGateway.BaiduApiKey = service.SettingService.GetSetting(model.SettingKeyBaiduApiKey)
	if baiduGateway.BaiduApiKey == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyBaiduApiKey)
	}
	secretKey := service.SettingService.GetSetting(model.SettingKeyBaiduSecretKey)
	if secretKey == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyBaiduSecretKey)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(secretKey), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	baiduGateway.BaiduSecretKey = string(decrypt)
	return baiduGateway, nil
}
func (gateway Baidu) Scope() string {
	return BaiduUserScope
}
func (gateway Baidu) GrantType() string {
	return BaiduGrantType
}

func (gateway Baidu) AuthorizeUrl(scope string, redirect string, state string) (string, string, error) {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", gateway.BaiduApiKey)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	query.Add("display", BaiduAuthorizeDisplay)
	query.Add("force_login", "1")
	query.Add("qrcode", "1")
	queryString := query.Encode()
	return BaiduAuthorizeUrl + "?" + queryString, state, nil
}

func (gateway Baidu) AccessToken(callbackData map[string]string, redirect string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(BaiduAccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"grant_type":    gateway.GrantType(),
			"code":          callbackData["code"],
			"client_id":     gateway.BaiduApiKey,
			"client_secret": gateway.BaiduSecretKey,
			"redirect_uri":  redirect,
		},
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
		"refreshToken": body.Get("refresh_token").String(),
	}, nil
}

func (gateway Baidu) User(accessToken string) (map[string]string, error) {
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
