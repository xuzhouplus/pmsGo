package baidu

import (
	"fmt"
	"github.com/idoubi/goz"
	"github.com/tidwall/gjson"
)

func PreCreate(accessToken string, path string, size string, isDir string, blockList []string) (*gjson.Result, error) {
	client := goz.NewClient()
	response, err := client.Get(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "precreate",
			"access_token": accessToken,
		},
		JSON: map[string]interface{}{
			"path":       path,
			"size":       size,
			"isdir":      isDir,
			"autoinit":   1,
			"rtype":      0,
			"block_list": blockList,
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
		return nil, fmt.Errorf("获取容量失败：%v", body.Get("errmsg").String())
	}
	return body, nil
}

func Upload(accessToken string, path string, uploadId string, partSeq int, file []byte) (*gjson.Result, error) {
	client := goz.NewClient()
	response, err := client.Post(UploadUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/x-www-form-urlencoded",
			"Accept":       "application/json",
		},
		Query: map[string]interface{}{
			"method":       "upload",
			"access_token": accessToken,
			"type":         "tmpfile",
			"path":         path,
			"uploadid":     partSeq,
		},
		FormParams: map[string]interface{}{
			"file": file,
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
		return nil, fmt.Errorf("上传文件失败：%v", body.Get("errmsg").String())
	}
	return body, nil
}

func Create(accessToken string, path string, size string, isDir string, blockList []string) (*gjson.Result, error) {
	client := goz.NewClient()
	response, err := client.Post(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "create",
			"access_token": accessToken,
		},
		JSON: map[string]interface{}{
			"path":       path,
			"size":       size,
			"isdir":      isDir,
			"autoinit":   1,
			"rtype":      0,
			"block_list": blockList,
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
		return nil, fmt.Errorf("创建文件失败：%v", body.Get("errmsg").String())
	}
	return body, nil
}
