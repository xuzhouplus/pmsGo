package baidu

import (
	"fmt"
	"github.com/idoubi/goz"
	"github.com/tidwall/gjson"
)

func FileManager(accessToken string, operate string, fileList []map[string]string) ([]gjson.Result, error) {
	client := goz.NewClient()
	response, err := client.Get(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]interface{}{
			"method":       "filemanager",
			"access_token": accessToken,
			"opera":        operate,
		},
		JSON: map[string]interface{}{
			"async":    0,
			"filelist": fileList,
			"ondup":    "fail",
		},
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	if body.Get("errno").Int() != 0 {
		return nil, fmt.Errorf("管理失败：%v", body.Get("errmsg").String())
	}

	return body.Get("list").Array(), nil
}