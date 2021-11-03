package baidu

import (
	"github.com/idoubi/goz"
	"log"
	"net/url"
)

const UserScope = "basic,netdisk"
const GrantType = "authorization_code"
const AuthorizeDisplay = "page"

const (
	AuthorizeUrl   = "https://openapi.baidu.com/oauth/2.0/authorize"
	AccessTokenUrl = "https://openapi.baidu.com/oauth/2.0/token"
	AccessUserUrl  = "https://openapi.baidu.com/rest/2.0/passport/users/getInf"
)

type AccessTokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func Authorize(scope string, redirect string, state string) (string, string, error) {
	if scope == "" {
		scope = UserScope
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", Baidu.BaiduApiKei)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	query.Add("display", AuthorizeDisplay)
	query.Add("force_login", "1")
	query.Add("qrcode", "1")
	queryString := query.Encode()
	return AuthorizeUrl + "?" + queryString, state, nil
}

func AccessToken(callbackData map[string]string, redirect string) (*AccessTokenResponse, error) {
	client := goz.NewClient()
	response, err := client.Get(AccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"grant_type":    GrantType,
			"code":          callbackData["code"],
			"client_id":     Baidu.BaiduApiKei,
			"client_secret": Baidu.BaiduSecretKey,
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
	accessTokenResponse := &AccessTokenResponse{}
	accessTokenResponse.AccessToken = body.Get("access_token").String()
	accessTokenResponse.RefreshToken = body.Get("refresh_token").String()
	accessTokenResponse.ExpiresIn = body.Get("expires_in").Int()
	accessTokenResponse.Scope = body.Get("scope").String()
	return accessTokenResponse, nil
}

func FreshToken(refreshToken string) (*AccessTokenResponse, error) {
	client := goz.NewClient()
	response, err := client.Get(AccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
			"client_id":     Baidu.BaiduApiKei,
			"client_secret": Baidu.BaiduSecretKey,
		},
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	accessTokenResponse := &AccessTokenResponse{}
	accessTokenResponse.AccessToken = body.Get("access_token").String()
	accessTokenResponse.RefreshToken = body.Get("refresh_token").String()
	accessTokenResponse.ExpiresIn = body.Get("expires_in").Int()
	accessTokenResponse.Scope = body.Get("scope").String()
	return accessTokenResponse, nil
}
func User(accessToken string) (map[string]string, error) {
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
	return map[string]string{
		"avatar":   body.Get("portrait").String(),
		"channel":  "0",
		"nickname": body.Get("username").String(),
		"gender":   body.Get("sex").String(),
		"open_id":  body.Get("openid").String(),
		"union_id": body.Get("unionid").String(),
	}, nil
}
