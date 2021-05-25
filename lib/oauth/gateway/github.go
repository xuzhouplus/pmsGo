package gateway

import (
	"errors"
	"github.com/idoubi/goz"
	"log"
	"net/url"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
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
	State        string `json:"state"`
}

type GitHubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func NewGitHub() (*GitHub, error) {
	gitHub := &GitHub{}
	settings := service.SettingService.GetSettings(model.GithubSettingModel.Keys())
	if settings == nil {
		return nil, errors.New("GitHub配置获取失败")
	}
	gitHub.GithubAppId = settings[model.SettingKeyGithubAppId]
	if gitHub.GithubAppId == "" {
		return nil, errors.New("缺少配置:" + model.SettingKeyGithubAppId)
	}
	gitHub.GithubApplicationName = settings[model.SettingKeyGithubApplicationName]
	if gitHub.GithubApplicationName == "" {
		return nil, errors.New("缺少配置:" + model.SettingKeyGithubApplicationName)
	}
	secret := settings[model.SettingKeyGithubAppSecret]
	if secret == "" {
		return nil, errors.New("缺少配置:" + model.SettingKeyGithubAppSecret)
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
func (gateway GitHub) AuthorizeUrl(scope string, redirect string, state string) (string, error) {
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
	return GitHubAuthorizeUrl + "?" + queryString, nil
}

func (gateway GitHub) AccessToken(code string, redirect string, state string) (string, error) {
	requestData := &GitHubAccessTokenRequest{gateway.GithubAppId, gateway.GithubAppSecret, code, redirect, state}
	client := goz.NewClient()
	response, err := client.Post(GitHubAccessTokenUrl, goz.Options{
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

func (gateway GitHub) User(accessToken string) (map[string]string, error) {
	client := goz.NewClient()
	response, err := client.Get(GitHubAccessTokenUrl, goz.Options{
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
	return nil, nil
}
