package gateway

import (
	"errors"
	"github.com/idoubi/goz"
	"net/url"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
)

const GiteeGatewayType = "Gitee"
const (
	GiteeAuthorizeUrl   = "https://gitee.com/oauth/authorize"
	GiteeAccessTokenUrl = "https://gitee.com/oauth/token"
	GiteeAccessUserUrl  = "https://gitee.com/api/v5/user"
)
const GiteeUserScope = "user_info"

const GiteeGrantType = "authorization_code"

type Gitee struct {
	GiteeApplicationName string
	GiteeAppId           string
	GiteeAppSecret       string
}

type GiteeAccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type GiteeAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func NewGitee() (*Gitee, error) {
	gitee := &Gitee{}
	settings := service.SettingService.GetSettings(model.GiteeSettingModel.Keys())
	if settings == nil {
		return nil, errors.New("Gitee配置获取失败")
	}
	gitee.GiteeAppId = settings[model.SettingKeyGiteeAppId]
	if gitee.GiteeAppId == "" {
		return nil, errors.New("缺少配置:" + model.SettingKeyGiteeAppId)
	}
	gitee.GiteeApplicationName = settings[model.SettingKeyGiteeApplicationName]
	if gitee.GiteeApplicationName == "" {
		return nil, errors.New("缺少配置:" + model.SettingKeyGiteeApplicationName)
	}
	secret := settings[model.SettingKeyGiteeAppSecret]
	if secret == "" {
		return nil, errors.New("缺少配置:" + model.SettingKeyGiteeAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(secret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gitee.GiteeAppSecret = string(decrypt)
	return gitee, nil
}
func (gateway Gitee) Scope() string {
	return GiteeUserScope
}
func (gateway Gitee) GrantType() string {
	return GiteeGrantType
}
func (gateway Gitee) AuthorizeUrl(scope string, redirect string, state string) (string, string, error) {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", gateway.GiteeAppId)
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	query.Add("response_type", "code")
	queryString := query.Encode()
	return GiteeAuthorizeUrl + "?" + queryString, state, nil
}

func (gateway Gitee) AccessToken(callbackData map[string]string, redirect string) (string, error) {
	requestData := &GiteeAccessTokenRequest{
		ClientId:     gateway.GiteeAppId,
		ClientSecret: gateway.GiteeAppSecret,
		Code:         callbackData["code"],
		RedirectUri:  redirect,
		GrantType:    gateway.GrantType(),
	}
	client := goz.NewClient()
	response, err := client.Post(GiteeAccessTokenUrl, goz.Options{
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

func (gateway Gitee) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(GiteeAccessUserUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"User-Agent":      gateway.GiteeApplicationName,
			"UserModel-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
			"Authorization":   "token " + accessToken,
		},
		Query: map[string]interface{}{
			"access_token": accessToken,
		},
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	if body.Get("root").Exists() {
		return map[string]string{
			"avatar":   body.Get("root.avatar_url").String(),
			"channel":  "0",
			"nickname": body.Get("root.name").String(),
			"gender":   "0",
			"open_id":  body.Get("root.id").String(),
			"union_id": body.Get("root.id").String(),
		}, nil
	}
	return map[string]string{
		"avatar":   body.Get("avatar_url").String(),
		"channel":  "0",
		"nickname": body.Get("name").String(),
		"gender":   "0",
		"open_id":  body.Get("id").String(),
		"union_id": body.Get("id").String(),
	}, nil
}
