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

const LineGatewayType = "line"

const LineScopeType = "profile%20openid"
const LineGrantType = "authorization_code"
const (
	LineAuthorizeUrl   = "https://access.line.me/oauth2/v2.1/authorize"
	LineAccessTokenUrl = "https://api.line.me/oauth2/v2.1/token"
	LineAccessUserUrl  = "https://api.line.me/v2/profile"
)

type LineAccessTokenRequest struct {
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

type Line struct {
	LineAppId     string
	LineAppSecret string
}

func NewLine() (*Line, error) {
	gateway := &Line{}
	appId := service.SettingService.GetSetting(model.SettingKeyLineAppId)
	if appId == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyLineAppId)
	}
	gateway.LineAppId = appId
	appSecret := service.SettingService.GetSetting(model.SettingKeyLineAppSecret)
	if appSecret == "" {
		return nil, fmt.Errorf("缺少配置：%v", model.SettingKeyLineAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appSecret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gateway.LineAppSecret = string(decrypt)
	return gateway, nil
}

func (gateway Line) Scope() string {
	return LineScopeType
}

func (gateway Line) GrantType() string {
	return LineGrantType
}

func (gateway Line) AuthorizeUrl(scope string, redirect string, state string) string {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", gateway.LineAppId)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	queryString := query.Encode()
	return LineAuthorizeUrl + "?" + queryString
}

func (gateway Line) AccessToken(code string, redirect string, state string) (string, error) {
	client := goz.NewClient()
	response, err := client.Post(LineAccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/x-www-form-urlencoded",
			"Accept":       "application/json",
		},
		FormParams: map[string]interface{}{
			"code":          code,
			"grant_type":    gateway.GrantType(),
			"client_id":     gateway.LineAppId,
			"client_secret": gateway.LineAppSecret,
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

func (gateway Line) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(LineAccessUserUrl, goz.Options{
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
	log.Println(body)
	return map[string]string{
		"avatar":   body.Get("pictureUrl").String() + "/large",
		"channel":  "0",
		"nickname": body.Get("displayName").String(),
		"gender":   "1",
		"open_id":  body.Get("userId").String(),
		"union_id": body.Get("userId").String(),
	}, nil
}
