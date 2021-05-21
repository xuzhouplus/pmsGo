package oauth

import (
	"fmt"
	"pmsGo/lib/oauth/gateway"
	"pmsGo/lib/oauth/user"
)

type Oauth struct {
	Type     string
	Instance gateway.Gateway
}

func NewOauth(gatewayType string) (*Oauth, error) {
	var gatewayInstance gateway.Gateway
	var err error
	switch gatewayType {
	case gateway.AlipayGatewayType:
		gatewayInstance, err = gateway.NewAlipay()
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
func (oauth Oauth) User(accessToken string) (*user.User, error) {
	auth, err := oauth.Instance.User(accessToken)
	if err != nil {
		return nil, err
	}
	authUser := &user.User{
		Type:     oauth.Type,
		Avatar:   auth["avatar"],
		Channel:  auth["channel"],
		Nickname: auth["nickname"],
		Gender:   auth["gender"],
		OpenId:   auth["open_id"],
		UnionId:  auth["union_id"],
	}
	return authUser, nil
}
