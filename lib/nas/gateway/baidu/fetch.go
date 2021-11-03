package baidu

import (
	"fmt"
	"github.com/idoubi/goz"
	"github.com/tidwall/gjson"
)

func UserInfo(accessToken string) (*UserInfoResponse, error) {
	client := goz.NewClient()
	response, err := client.Get(NasUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "uinfo",
			"access_token": accessToken,
		},
	})
	if err != nil {
		return nil, err
	}
	body, err := response.GetParsedBody()
	if err != nil {
		return nil, err
	}
	userInfoResponse := &UserInfoResponse{}
	userInfoResponse.BaiduName = body.Get("baidu_name").String()
	userInfoResponse.NetdiskName = body.Get("netdisk_name").String()
	userInfoResponse.AvatarUrl = body.Get("avatar_url").String()
	userInfoResponse.VipType = body.Get("vip_type").Int()
	userInfoResponse.Uk = body.Get("uk").Int()
	return userInfoResponse, nil
}

func Quota(accessToken string) (*QuotaResponse, error) {
	client := goz.NewClient()
	response, err := client.Get(QuotaUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "uinfo",
			"access_token": accessToken,
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
	quotaResponse := &QuotaResponse{}
	quotaResponse.Total = body.Get("total").Int()
	quotaResponse.Free = body.Get("free").Int()
	quotaResponse.Expire = body.Get("expire").Bool()
	quotaResponse.Used = body.Get("used").Int()
	return quotaResponse, nil
}

func List(accessToken string, dir string, page int, limit int, order string, desc string) ([]gjson.Result, error) {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * limit
	client := goz.NewClient()
	response, err := client.Get(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "list",
			"access_token": accessToken,
			"dir":          dir,
			"start":        string(start),
			"limit":        string(limit),
			"order":        order,
			"desc":         desc,
			"web":          "web",
			"folder":       "0",
			"showempty":    "1",
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
		return nil, fmt.Errorf("获取文件列表失败：%v", body.Get("errmsg").String())
	}

	return body.Get("list").Array(), nil
}

// DocList https://pan.baidu.com/union/doc/Eksg0saqp
func DocList(accessToken string, dir string, page int, limit int, order string, desc string) ([]gjson.Result, error) {
	if page < 1 {
		page = 1
	}
	client := goz.NewClient()
	response, err := client.Get(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "doclist",
			"access_token": accessToken,
			"page":         string(page),
			"num":          string(limit),
			"order":        order,
			"desc":         desc,
			"web":          string(1),
			"parent_path":  dir,
			"recursion":    string(0),
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
		return nil, fmt.Errorf("获取文档列表失败：%v", body.Get("errmsg").String())
	}

	return body.Get("list").Array(), nil
}

// ImageList https://pan.baidu.com/union/doc/bksg0sayv
func ImageList(accessToken string, dir string, page int, limit int, order string, desc string) ([]gjson.Result, error) {
	if page < 1 {
		page = 1
	}
	client := goz.NewClient()
	response, err := client.Get(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "imagelist",
			"access_token": accessToken,
			"page":         string(page),
			"num":          string(limit),
			"order":        order,
			"desc":         desc,
			"web":          string(1),
			"parent_path":  dir,
			"recursion":    string(0),
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
		return nil, fmt.Errorf("获取图片列表失败：%v", body.Get("errmsg").String())
	}

	return body.Get("list").Array(), nil
}

// VideoList https://pan.baidu.com/union/doc/Sksg0saw0
func VideoList(accessToken string, dir string, page int, limit int, order string, desc string) ([]gjson.Result, error) {
	if page < 1 {
		page = 1
	}
	client := goz.NewClient()
	response, err := client.Get(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "videolist",
			"access_token": accessToken,
			"page":         string(page),
			"num":          string(limit),
			"order":        order,
			"desc":         desc,
			"web":          string(1),
			"parent_path":  dir,
			"recursion":    string(0),
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
		return nil, fmt.Errorf("获取视频列表失败：%v", body.Get("errmsg").String())
	}

	return body.Get("list").Array(), nil
}

// Search https://pan.baidu.com/union/doc/zksg0sb9z
func Search(accessToken string, dir string, page int, limit int, key string) ([]gjson.Result, error) {
	if page < 1 {
		page = 1
	}
	client := goz.NewClient()
	response, err := client.Get(FileUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]string{
			"method":       "search",
			"access_token": accessToken,
			"page":         string(page),
			"num":          string(limit),
			"key":          key,
			"web":          string(1),
			"dir":          dir,
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
		return nil, fmt.Errorf("搜索文件失败：%v", body.Get("errmsg").String())
	}

	return body.Get("list").Array(), nil
}

// FileMetas https://pan.baidu.com/union/doc/Fksg0sbcm
func FileMetas(accessToken string, dir string, fileIds []uint64) ([]gjson.Result, error) {
	client := goz.NewClient()
	response, err := client.Get(MultiMediaUrl, goz.Options{
		Headers: map[string]interface{}{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Query: map[string]interface{}{
			"method":       "filemetas",
			"access_token": accessToken,
			"fsids":        fileIds,
			"thumb":        1,
			"dlink":        0,
			"extra":        1,
			"path":         dir,
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
		return nil, fmt.Errorf("查询文件信息失败：%v", body.Get("errmsg").String())
	}

	return body.Get("list").Array(), nil
}


