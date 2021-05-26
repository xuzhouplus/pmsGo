package gateway

import (
	"crypto"
	"fmt"
	"github.com/idoubi/goz"
	"net/url"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
	"pmsGo/lib/security/rsa"
	"sort"
	"strings"
	"time"
)

const AlipayGatewayType = "alipay"
const (
	AlipayAuthorizeUrl = "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm"
	AlipayGatewayUrl   = "https://openapi.alipay.com/gateway.do"
)

const AlipayUserScope = "auth_user"
const AlipayGrantType = "authorization_code"

type AlipayAccessTokenRequest struct {
	AppId     string `json:"app_id"`
	Method    string `json:"method"`
	Charset   string `json:"charset"`
	SignType  string `json:"sign_type"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
	Format    string `json:"format"`
	Sign      string `json:"sign"`
}

type AlipayAccessUserRequest struct {
	AppId     string `json:"app_id"`
	Method    string `json:"method"`
	Charset   string `json:"charset"`
	SignType  string `json:"sign_type"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	AuthToken string `json:"auth_token"`
	Format    string `json:"format"`
	Sign      string `json:"sign"`
}

// Alipay @link https://opendocs.alipay.com/open/263/105808/**
type Alipay struct {
	AlipayAppId         string
	AlipayAppPrimaryKey []byte
	AlipayPublicKey     []byte
}

func NewAlipay() (*Alipay, error) {
	alipay := &Alipay{}
	alipay.AlipayAppId = service.SettingService.GetSetting(model.SettingKeyAlipayAppId)
	if alipay.AlipayAppId == "" {
		return nil, fmt.Errorf("缺少参数：%v", model.SettingKeyAlipayAppId)
	}
	publicKey := service.SettingService.GetSetting(model.SettingKeyAlipayPublicKay)
	if publicKey == "" {
		return nil, fmt.Errorf("缺少参数：%v", model.SettingKeyAlipayPublicKay)
	}
	alipay.AlipayPublicKey = rsa.FormatPublicKey(publicKey)
	appPrimaryKey := service.SettingService.GetSetting(model.SettingKeyAlipayAppPrimaryKey)
	if appPrimaryKey == "" {
		return nil, fmt.Errorf("缺少参数：%v", model.SettingKeyAlipayAppPrimaryKey)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appPrimaryKey), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	alipay.AlipayAppPrimaryKey = rsa.FormatPKCS8PrivateKey(string(decrypt))
	return alipay, nil
}
func (gateway Alipay) Scope() string {
	return AlipayUserScope
}

func (gateway Alipay) GrantType() string {
	return AlipayGrantType
}

func (gateway Alipay) signature(params map[string]interface{}) string {
	var pList = make([]string, 0, 0)

	for field, val := range params {
		if val != nil {
			pList = append(pList, field+"="+val.(string))
		}
	}
	sort.Strings(pList)
	src := strings.Join(pList, "&")
	signature, err := rsa.RSASignWithPKCS8([]byte(src), gateway.AlipayAppPrimaryKey, crypto.SHA256)
	if err != nil {
		panic(err.Error())
	}
	return base64.Encode(signature)
}

func (gateway Alipay) AuthorizeUrl(scope string, redirect string, state string) (string, string, error) {
	if scope == "" {
		scope = gateway.Scope()
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("app_id", gateway.AlipayAppId)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	queryString := query.Encode()
	return AlipayAuthorizeUrl + "?" + queryString, state, nil
}

func (gateway Alipay) AccessToken(callbackData map[string]string, redirect string) (string, error) {
	queryData := map[string]interface{}{
		"app_id":     gateway.AlipayAppId,
		"method":     "alipay.system.oauth.token",
		"charset":    "utf-8",
		"sign_type":  "RSA2",
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
		"version":    "1.0",
		"grant_type": gateway.GrantType(),
		"code":       callbackData["auth_code"],
		"format":     "JSON",
	}
	queryData["sign"] = gateway.signature(queryData)
	client := goz.NewClient()
	response, err := client.Post(AlipayGatewayUrl, goz.Options{
		Debug: true,
		Headers: map[string]interface{}{
			"Accept": "application/json",
		},
		FormParams: queryData,
	})
	if err != nil {
		return "", err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return "", err
	}
	responseData := body.Get("alipay_system_oauth_token_response")
	if responseData.Exists() {
		return responseData.Get("access_token").String(), nil
	}
	return "", fmt.Errorf("获取支付宝 ACCESS_TOKEN 出错：%v", body.String())
}

func (gateway Alipay) User(accessToken string) (map[string]string, error) {
	queryData := map[string]interface{}{
		"app_id":     gateway.AlipayAppId,
		"method":     "alipay.user.info.share",
		"charset":    "utf-8",
		"sign_type":  "RSA2",
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
		"version":    "1.0",
		"auth_token": accessToken,
		"format":     "JSON",
	}
	queryData["sign"] = gateway.signature(queryData)

	client := goz.NewClient()
	response, err := client.Post(AlipayGatewayUrl, goz.Options{
		Headers: map[string]interface{}{
			"Accept": "application/json",
		},
		FormParams: queryData,
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	responseData := body.Get("alipay_user_info_share_response")
	if responseData.Exists() {
		sex := "2"
		gender := responseData.Get("gender").String()
		if gender == "m" {
			sex = "1"
		}
		return map[string]string{
			"avatar":   responseData.Get("avatar").String(),
			"channel":  "0",
			"nickname": responseData.Get("nick_name").String(),
			"gender":   sex,
			"open_id":  responseData.Get("user_id").String(),
			"union_id": responseData.Get("user_id").String(),
		}, nil
	}
	return nil, fmt.Errorf("获取支付宝用户出错：%v", body.String())
}
