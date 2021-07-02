package helper

import (
	"net/url"
	"strings"
)

func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

func RawUrlEncode(str string) string {
	return strings.Replace(UrlEncode(str), "+", "%20", -1)
}
