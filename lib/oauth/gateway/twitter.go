package gateway

import (
	"crypto/sha1"
	"fmt"
	"github.com/idoubi/goz"
	"log"
	"net/url"
	"pmsGo/lib/config"
	"pmsGo/lib/helper"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
	"pmsGo/lib/security/random"
	model2 "pmsGo/model"
	service2 "pmsGo/service"
	"sort"
	"strings"
	"time"
)

const TwitterGatewayType = "twitter"

const TwitterScopeType = ""
const TwitterGrantType = ""
const (
	TwitterRequestTokenUrl = "https://api.twitter.com/oauth/request_token"
	TwitterAuthorizeUrl    = "https://api.twitter.com/oauth/authenticate"
	TwitterAccessTokenUrl  = "https://api.twitter.com/oauth/access_token"
	TwitterAccessUserUrl   = "https://api.twitter.com/1.1/account/verify_credentials.json"
)

type TwitterAccessTokenRequest struct {
}

type Twitter struct {
	TwitterAppId       string
	TwitterAppSecret   string
	TwitterTokenSecret string
}

func NewTwitter() (*Twitter, error) {
	gateway := &Twitter{}
	appId := service2.SettingService.GetSetting(model2.SettingKeyTwitterAppId)
	if appId == "" {
		return nil, fmt.Errorf("缺少配置：%v", model2.SettingKeyTwitterAppId)
	}
	gateway.TwitterAppId = appId
	appSecret := service2.SettingService.GetSetting(model2.SettingKeyTwitterAppSecret)
	if appSecret == "" {
		return nil, fmt.Errorf("缺少配置：%v", model2.SettingKeyTwitterAppSecret)
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(appSecret), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		return nil, err
	}
	gateway.TwitterAppSecret = string(decrypt)
	return gateway, nil
}

func (gateway Twitter) Scope() string {
	return TwitterScopeType
}

func (gateway Twitter) GrantType() string {
	return TwitterGrantType
}
func (gateway Twitter) Sign(method, uri string, params map[string]interface{}) (map[string]interface{}, map[string]interface{}) {
	defaultParams := map[string]interface{}{
		"oauth_consumer_key":     gateway.TwitterAppId,
		"oauth_nonce":            random.Uuid(false),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        time.Now().String(),
		"oauth_token":            "",
		"oauth_version":          "1.0",
	}
	for key, value := range params {
		defaultParams[key] = value.(string)
	}
	pList := make([]string, 0, 0)
	for field, val := range defaultParams {
		valStr := val.(string)
		valStr = helper.RawUrlEncode(valStr)
		pList = append(pList, field+"="+valStr)
	}
	sort.Strings(pList)
	signStr := strings.Join(pList, "&")
	signStr = strings.ToUpper(method) + "&" + helper.RawUrlEncode(uri) + "&" + helper.RawUrlEncode(signStr)
	signKey := gateway.TwitterAppSecret + "&" + gateway.TwitterTokenSecret
	hashEncrypted := encrypt.HashHmac(sha1.New, []byte(signStr), []byte(signKey), false)
	base64Encrypted := base64.Encode(hashEncrypted)
	authSign := helper.RawUrlEncode(base64Encrypted)
	defaultParams["oauth_signature"] = authSign
	authStr := "OAuth "
	for key, val := range defaultParams {
		authStr = authStr + key + "=\"" + val.(string) + "\", "
	}
	headers := map[string]interface{}{
		"Authorization": strings.TrimRight(authStr, ", "),
	}
	return defaultParams, headers
}
func (gateway Twitter) RequestToken(redirect string) (string, error) {
	formParams := map[string]interface{}{
		"oauth_callback": redirect,
	}
	oauthParams, headers := gateway.Sign("POST", TwitterRequestTokenUrl, formParams)
	client := goz.NewClient()
	response, err := client.Post(TwitterRequestTokenUrl, goz.Options{
		Headers:    headers,
		FormParams: oauthParams,
	})
	if err != nil {
		return "", err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return "", err
	}
	return body.Get("oauth_token").String(), nil
}

func (gateway Twitter) AuthorizeUrl(scope string, redirect string, state string) (string, string, error) {
	token, err := gateway.RequestToken(redirect)
	if err != nil {
		return "", "", err
	}
	uri := url.URL{}
	query := uri.Query()
	query.Add("oauth_token", token)
	queryString := query.Encode()
	return TwitterAuthorizeUrl + "?" + queryString, token, nil
}

// AccessToken 授权回调，参数为oauth_token和oauth_verifier
func (gateway *Twitter) AccessToken(callbackData map[string]string, redirect string) (string, error) {
	formParams := map[string]interface{}{
		"oauth_token":    callbackData["oauth_token"],
		"oauth_verifier": callbackData["oauth_verifier"],
	}
	oauthParams, headers := gateway.Sign("POST", TwitterAccessTokenUrl, formParams)
	client := goz.NewClient()
	response, err := client.Post(TwitterAccessTokenUrl, goz.Options{
		Headers:    headers,
		FormParams: oauthParams,
	})
	if err != nil {
		return "", err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return "", err
	}
	gateway.TwitterTokenSecret = body.Get("oauth_token_secret").String()
	return body.Get("oauth_token").String(), nil
}

func (gateway Twitter) User(accessToken string) (map[string]string, error) {
	formParams := map[string]interface{}{
		"oauth_token":        accessToken,
		"oauth_token_secret": gateway.TwitterTokenSecret,
	}
	oauthParams, headers := gateway.Sign("GET", TwitterAccessUserUrl, formParams)
	client := goz.NewClient()
	response, err := client.Get(TwitterAccessUserUrl, goz.Options{
		Headers: headers,
		Query:   oauthParams,
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
		"avatar":   body.Get("profile_image_url_https").String(),
		"channel":  "0",
		"nickname": body.Get("name").String(),
		"gender":   "2",
		"open_id":  body.Get("id_str").String(),
		"union_id": body.Get("id_str").String(),
	}, nil
}
