package baidu

import (
	"github.com/idoubi/goz"
	"log"
	"net/url"
	"pmsGo/lib/nas/gateway"
)

const UserScope = "basic,netdisk"
const GrantType = "authorization_code"
const AuthorizeDisplay = "page"

const (
	AuthorizeUrl   = "https://openapi.baidu.com/oauth/2.0/authorize"
	AccessTokenUrl = "https://openapi.baidu.com/oauth/2.0/token"
	AccessUserUrl  = "https://openapi.baidu.com/rest/2.0/passport/users/getInf"
)

func (service baidu) Authorize(redirect string, state string) (string, error) {
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", service.BaiduApiKey)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirect)
	query.Add("scope", UserScope)
	query.Add("state", state)
	query.Add("display", AuthorizeDisplay)
	query.Add("force_login", "1")
	query.Add("qrcode", "1")
	queryString := query.Encode()
	return AuthorizeUrl + "?" + queryString, nil
}

func (service baidu) AccessToken(callbackData map[string]string, redirect string) (*gateway.AccessToken, error) {
	client := goz.NewClient()
	response, err := client.Get(AccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"grant_type":    GrantType,
			"code":          callbackData["code"],
			"client_id":     service.BaiduApiKey,
			"client_secret": service.BaiduSecretKey,
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
	accessTokenResponse := &gateway.AccessToken{}
	accessTokenResponse.AccessToken = body.Get("access_token").String()
	accessTokenResponse.RefreshToken = body.Get("refresh_token").String()
	accessTokenResponse.ExpiresIn = body.Get("expires_in").Int()
	accessTokenResponse.Scope = body.Get("scope").String()
	return accessTokenResponse, nil
}

func (service baidu) FreshToken(refreshToken string) (*gateway.AccessToken, error) {
	client := goz.NewClient()
	response, err := client.Get(AccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
			"client_id":     service.BaiduApiKey,
			"client_secret": service.BaiduSecretKey,
		},
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	accessTokenResponse := &gateway.AccessToken{}
	accessTokenResponse.AccessToken = body.Get("access_token").String()
	accessTokenResponse.RefreshToken = body.Get("refresh_token").String()
	accessTokenResponse.ExpiresIn = body.Get("expires_in").Int()
	accessTokenResponse.Scope = body.Get("scope").String()
	return accessTokenResponse, nil
}

func (service baidu) User(accessToken string) (*gateway.User, error) {
	client := goz.NewClient()
	response, err := client.Get(AccessUserUrl, goz.Options{
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
	return &gateway.User{
		OpenId:   body.Get("openid").String(),
		UnionId:  body.Get("unionid").String(),
		Channel:  "0",
		Nickname: body.Get("username").String(),
		Gender:   body.Get("sex").String(),
		Avatar:   body.Get("portrait").String(),
		Type:     GatewayType,
	}, nil
}
