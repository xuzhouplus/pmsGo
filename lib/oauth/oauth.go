package oauth

import (
	"fmt"
	"pmsGo/lib/oauth/gateway"
)

type Oauth struct {
	Type     string
	Instance gateway.Gateway
}

type User struct {
	OpenId   string `json:"open_id"`  //用户唯一id
	UnionId  string `json:"union_id"` //微信union_id
	Channel  string `json:"channel"`  //登录类型请查看 \\tinymeng\\OAuth2\\Helper\\ConstCode
	Nickname string `json:"nickname"` //昵称
	Gender   string `json:"gender"`   //0=>未知 1=>男 2=>女   twitter和line不会返回性别，所以这里是0，Facebook根据你的权限，可能也不会返回，所以也可能是0
	Avatar   string `json:"avatar"`   //头像
	Type     string `json:"type"`     //授权类型
}

func NewOauth(gatewayType string) (*Oauth, error) {
	var gatewayInstance gateway.Gateway
	var err error
	switch gatewayType {
	case gateway.GitHubGatewayType:
		gatewayInstance, err = gateway.NewGitHub()
	case gateway.BaiduGatewayType:
		gatewayInstance, err = gateway.NewBaidu()
	default:
		return nil, fmt.Errorf("不支持的类型：%v", gatewayType)
	}
	if err != nil {
		return nil, err
	}
	oauth := &Oauth{Type: gatewayType, Instance: gatewayInstance}
	return oauth, nil
}
func (oauth Oauth) AuthorizeUrl(scope string, redirect string, state string) string {
	return oauth.Instance.AuthorizeUrl(scope, redirect, state)
}
func (oauth Oauth) AccessToken(code string, redirect string, state string) (string, error) {
	token, err := oauth.Instance.AccessToken(code, redirect, state)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (oauth Oauth) User(accessToken string) (*User, error) {
	user, err := oauth.Instance.User(accessToken)
	if err != nil {
		return nil, err
	}
	authUser := &User{
		Type:     oauth.Type,
		Avatar:   user["avatar"],
		Channel:  user["channel"],
		Nickname: user["nickname"],
		Gender:   user["gender"],
		OpenId:   user["open_id"],
		UnionId:  user["union_id"],
	}
	return authUser, nil
}
