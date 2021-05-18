package image

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"pmsGo/lib/security"
	"regexp"
	"time"
)

type Instance struct {
	ctx       *gin.Context
	File      string `json:"image"`
	MimeType  string `json:"mimeType"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

func Upload(ctx *gin.Context, fieldName string, subDir string) (*Instance, error) {
	helper := &Instance{}
	helper.ctx = ctx
	file, err := ctx.FormFile(fieldName)
	if err != nil {
		return nil, err
	}
	log.Printf("Upload image:%v \n", file.Header)
	helper.MimeType = file.Header.Get("Content-Type")
	err = validateMimeType(helper.MimeType)
	if err != nil {
		return nil, err
	}
	helper.Extension = path.Ext(file.Filename)
	err = validateExtensions(helper.Extension)
	if err != nil {
		return nil, err
	}
	helper.Size = file.Size
	err = validateMaxSize(helper.Size)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	date := now.Format("2006-01-02")
	guid := security.Uuid(false)
	helper.File = subDir + "/" + date + "/" + guid + helper.Extension
	filePath := Settings.path + helper.File
	filePath = filepath.ToSlash(filePath)
	err = mkdir(path.Dir(filePath))
	if err != nil {
		return nil, err
	}
	err = ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		return nil, err
	}
	return helper, nil
}

func Base64Upload(ctx *gin.Context, fieldName string, subDir string) (*Instance, error) {
	helper := &Instance{}
	helper.ctx = ctx
	//获取请求数据
	postData := make(map[string]interface{})
	ctx.ShouldBind(&postData)
	if postData[fieldName] == nil {
		return nil, errors.New("字段数据为空")
	}
	//获取请求内容
	base64Content := postData[fieldName].(string)
	//验证图片格式
	matched, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64Content)
	if !matched {
		return nil, errors.New("数据格式错误")
	}
	//解析图片mimeType、extension和内容
	reg, _ := regexp.Compile(`^data:(image\/(\w+));base64,`)
	allData := reg.FindAllSubmatch([]byte(base64Content), 2)
	helper.MimeType = string(allData[0][1])        //image/png
	helper.Extension = "." + string(allData[0][2]) //png ，jpeg 后缀获取
	base64Str := reg.ReplaceAllString(base64Content, "")
	//拼接图片保存路径
	now := time.Now()
	date := now.Format("2006-01-02")
	guid := security.Uuid(false)
	helper.File = subDir + "/" + date + "/" + guid + helper.Extension
	filePath := Settings.path + helper.File
	filePath = filepath.ToSlash(filePath)
	//创建文件夹
	err := mkdir(path.Dir(filePath))
	if err != nil {
		return nil, err
	}
	//写入文件
	err = ioutil.WriteFile(filePath, []byte(base64Str), os.ModePerm)
	if err != nil {
		return nil, err
	}
	//获取文件大小
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	helper.Size = stat.Size()
	return helper, nil
}

func (helper Instance) Path() Path {
	return Path(filepath.FromSlash(Settings.path + helper.File))
}

func (helper Instance) Url() Url {
	return Url(filepath.ToSlash(Settings.url + helper.File))
}
