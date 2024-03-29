package gateway

import (
	"errors"
	"github.com/idoubi/goz"
	"log"
	"net/url"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
	model2 "pmsGo/model"
	service2 "pmsGo/service"
)

const GitHubGatewayType = "github"
const (
	GitHubAuthorizeUrl   = "https://github.com/login/oauth/authorize"
	GitHubAccessTokenUrl = "https://github.com/login/oauth/access_token"
	GitHubAccessUserUrl  = "https://api.github.com/user"
)
const GitHubUserScope = "user"

const GitHubGrantType = "authorization_code"

type GitHub struct {
	GithubApplicationName string
	GithubAppId           string
	GithubAppSecret       string
}

type GitHubAccessTokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type GitHubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func NewGitHub() (*GitHub, error) {
	gitHub := &GitHub{}
	settings := service2.SettingService.GetSettings(model2.GithubSettingModel.Keys())
	if settings == nil {
		return nil, errors.New("GitHub配置获取失败")
	}
	gitHub.GithubAppId = settings[model2.SettingKeyGithubAppId]
	if gitHub.GithubAppId == "" {
		return nil, errors.New("缺少配置:" + model2.SettingKeyGithubAppId)
	}
	gitHub.GithubApplicationName = settings[model2.SettingKeyGithubApplicationName]
	if gitHub.GithubApplicationName == "" {
		return nil, errors.New("缺少配置:" + model2.SettingKeyGithubApplicationName)
	}
	secret := settings[model2.SettingKeyGithubAppSecret]
	if secret == "" {
		return nil, errors.New("缺少配置:" + model2.SettingKeyGithubAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(secret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gitHub.GithubAppSecret = string(decrypt)
	return gitHub, nil
}
func (gateway GitHub) Scope() string {
	return GitHubUserScope
}
func (gateway GitHub) GrantType() string {
	return GitHubGrantType
}
func (gateway GitHub) AuthorizeUrl(scope string, redirect string, state string) (string, string, error) {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("client_id", gateway.GithubAppId)
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	query.Add("allow_signup", "true")
	queryString := query.Encode()
	return GitHubAuthorizeUrl + "?" + queryString, state, nil
}

func (gateway GitHub) AccessToken(callbackData map[string]string, redirect string) (map[string]string, error) {
	requestData := &GitHubAccessTokenRequest{
		ClientId:     gateway.GithubAppId,
		ClientSecret: gateway.GithubAppSecret,
		Code:         callbackData["code"],
		RedirectUri:  redirect,
		GrantType:    gateway.GrantType(),
	}
	client := goz.NewClient()
	response, err := client.Post(GitHubAccessTokenUrl, goz.Options{
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
		"refreshToken": "",
	}, nil
}

func (gateway GitHub) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(GitHubAccessUserUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type":    "application/json",
			"Accept":          "application/json",
			"User-Agent":      gateway.GithubApplicationName,
			"UserModel-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
			"Authorization":   "token " + accessToken,
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
		"avatar":   body.Get("avatar_url").String(),
		"channel":  "0",
		"nickname": body.Get("name").String(),
		"gender":   "0",
		"open_id":  body.Get("id").String(),
		"union_id": body.Get("id").String(),
	}, nil
}
