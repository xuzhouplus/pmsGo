package baidu

import (
	"fmt"
	"pmsGo/lib/config"
	"pmsGo/lib/log"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/encrypt"
	"pmsGo/model"
	"pmsGo/service"
)

const MultiMediaUrl = "https://pan.baidu.com/rest/2.0/xpan/multimedia"

const FileUrl = "https://pan.baidu.com/rest/2.0/xpan/file"

const NasUrl = "https://pan.baidu.com/rest/2.0/xpan/nas"

const QuotaUrl = "https://pan.baidu.com/api/quota"

const UploadUrl = "https://d.pcs.baidu.com/rest/2.0/pcs/superfile2"

type UserInfoResponse struct {
	BaiduName   string //百度账号
	NetdiskName string //网盘账号
	AvatarUrl   string //头像地址
	VipType     int64  //会员类型，0普通用户、1普通会员、2超级会员
	Uk          int64  //用户ID
}

type FileInfo struct {
	FsId           uint64    //文件在云端的唯一标识ID
	Path           string    //文件的绝对路径
	ServerFileName string    //文件名称
	Size           uint      //文件大小，单位B
	ServerMtime    uint      //文件在服务器修改时间
	ServerCtime    uint      //文件在服务器创建时间
	LocalMtime     uint      //文件在客户端修改时间
	LocalCtime     uint      //文件在客户端创建时间
	IsDir          uint      //是否目录，0 文件、1 目录
	Category       uint      //文件类型，1 视频、2 音频、3 图片、4 文档、5 应用、6 其他、7 种子
	Md5            string    //文件的md5值，只有是文件类型时，该KEY才存在
	DirEmpty       int       //该目录是否存在子目录， 只有请求参数带WEB且该条目为目录时，该KEY才存在， 0为存在， 1为不存在
	Thumbs         [3]string //只有请求参数带WEB且该条目分类为图片时，该KEY才存在，包含三个尺寸的缩略图URL
}

type QuotaResponse struct {
	Total  int64
	Free   int64
	Expire bool
	Used   int64
}

type baidu struct {
	BaiduAppName    string
	BaiduApiKey    string
	BaiduSecretKey string
}

var Baidu *baidu

func init() {
	Baidu := &baidu{}
	Baidu.BaiduApiKey = service.SettingService.GetSetting(model.SettingKeyBaiduApiKey)
	if Baidu.BaiduApiKey == "" {
		log.Error(fmt.Errorf("缺少配置：%v", model.SettingKeyBaiduApiKey))
		return
	}
	secretKey := service.SettingService.GetSetting(model.SettingKeyBaiduSecretKey)
	if secretKey == "" {
		log.Error(fmt.Errorf("缺少配置：%v", model.SettingKeyBaiduSecretKey))
		return
	}
	decrypt, err := encrypt.Decrypt(base64.Decode(secretKey), []byte(config.Config.Web.Security["salt"]))
	if err != nil {
		log.Error(err)
	}
	Baidu.BaiduSecretKey = string(decrypt)
}
