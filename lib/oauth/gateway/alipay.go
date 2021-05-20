package gateway

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/idoubi/goz"
	"log"
	"net/url"
	"pmsGo/app/model"
	"pmsGo/app/service"
	"pmsGo/lib/config"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
	"sort"
	"strings"
	"time"
)

const (
	AlipayAuthorizeUrl   = "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm"
	AlipayAccessTokenUrl = "https://openapi.alipay.com/gateway.do"
)

const AlipayUserScope = "auth_user"
const AlipayGrantType = "authorization_code"

// Alipay @link https://opendocs.alipay.com/open/263/105808/**
type Alipay struct {
	AlipayAppId         string
	AlipayAppPrimaryKey string
	AlipayPublicKey     string
}

func NewAlipay() (*Alipay, error) {
	alipay := &Alipay{}
	alipay.AlipayAppId = service.SettingService.GetSetting(model.SettingKeyAlipayAppId)
	if alipay.AlipayAppId == "" {
		return nil, fmt.Errorf("缺少参数：%v", model.SettingKeyAlipayAppId)
	}
	alipay.AlipayPublicKey = service.SettingService.GetSetting(model.SettingKeyAlipayPublicKay)
	if alipay.AlipayPublicKey == "" {
		return nil, fmt.Errorf("缺少参数：%v", model.SettingKeyAlipayPublicKay)
	}
	appPrimaryKey := service.SettingService.GetSetting(model.SettingKeyAlipayAppPrimaryKey)
	if appPrimaryKey == "" {
		return nil, fmt.Errorf("缺少参数：%v", model.SettingKeyAlipayAppPrimaryKey)
	}
	decrypt, err := encrypt.Decrypt([]byte(appPrimaryKey), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	alipay.AlipayAppPrimaryKey = string(decrypt)
	return alipay, nil
}
func (gateway Alipay) Scope() string {
	return AlipayUserScope
}

func (gateway Alipay) GrantType() string {
	return AlipayGrantType
}

func (gateway Alipay) signature(params map[string]string) string {
	var pList = make([]string, 0, 0)

	for field, val := range params {
		pList = append(pList, field+"="+val)
	}
	sort.Strings(pList)
	var src = strings.Join(pList, "&")
	h := sha256.New()
	h.Write([]byte(src))
	hashed := h.Sum(nil)
	block, _ := pem.Decode([]byte(gateway.AlipayAppPrimaryKey))
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return base64.Encode(signature)
}

func (gateway Alipay) AuthorizeUrl(scope string, redirect string, state string) string {
	if scope == "" {
		scope = gateway.Scope()
	}
	url := url.URL{}
	query := url.Query()
	query.Add("app_id", gateway.AlipayAppId)
	query.Add("response_type", "code")
	query.Add("redirect_uri", redirect)
	query.Add("scope", scope)
	query.Add("state", state)
	queryString := query.Encode()
	return BaiduAuthorizeUrl + "?" + queryString
}

func (gateway Alipay) AccessToken(code string, redirect string, state string) (string, error) {
	queryData := map[string]string{
		"app_id":     gateway.AlipayAppId,
		"method":     "alipay.system.oauth.token",
		"charset":    "utf-8",
		"sign_type":  "RSA2",
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
		"version":    "1.0",
		"grant_type": gateway.GrantType(),
		"code":       code,
	}
	queryData["sign"] = gateway.signature(queryData)

	client := goz.NewClient()
	response, err := client.Post(AlipayAccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		JSON: queryData,
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
	queryData := map[string]string{
		"app_id":     gateway.AlipayAppId,
		"method":     "alipay.user.info.share",
		"charset":    "utf-8",
		"sign_type":  "RSA2",
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
		"version":    "1.0",
		"auth_token": accessToken,
	}
	queryData["sign"] = gateway.signature(queryData)

	client := goz.NewClient()
	response, err := client.Post(AlipayAccessTokenUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		JSON: queryData,
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	responseData := body.Get("alipay_user_info_share_response")
	log.Println(responseData)
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
