package baidu

import (
	"fmt"
	"github.com/idoubi/goz"
)

func Streaming(accessToken string, path string, streamType string, adToken string) (map[string]interface{}, error) {
	client := goz.NewClient()
	response, err := client.Get(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "streaming",
			"access_token": accessToken,
			"path":         path,
			"type":         streamType,
			"adToken":      adToken,
		},
	})
	if err != nil {
		return nil, err
	}
	if adToken != "" {
		body, err := response.GetParsedBody()
		if err != nil {
			return nil, err
		}
		if body.Get("errno").Int() != 0 {
			return nil, fmt.Errorf("创建文件失败：%v", body.Get("errmsg").String())
		}
		return map[string]interface{}{
			"addTime": body.Get("adTime").Int(),
			"adToken": body.Get("adToken").String(),
			"ltime":   body.Get("ltime").Int(),
		}, nil
	}
	return map[string]interface{}{
		"m3u8": response.GetBody(),
	}, nil
}
