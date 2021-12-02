package file

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"pmsGo/lib/config"
	"pmsGo/lib/log"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/random"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	TypeVideo   = "video"
	TypeImage   = "image"
	TypeUnknown = "unknown"
)

const SubDir = "baseOnFileType"

type Url string

type Path string

func init() {
	log.Debugf("Upload: %v \n", config.Config.Web.Upload)
}

// FormUpload 通用表单文件上传
func FormUpload(ctx *gin.Context, fieldName string, subDir string) (*Upload, error) {
	upload := &Upload{}
	upload.ctx = ctx
	file, err := ctx.FormFile(fieldName)
	if err != nil {
		return nil, err
	}
	log.Debugf("Upload file:%v \n", file.Header)
	upload.MimeType = file.Header.Get("Content-Type")
	upload.FileType = GetFileType(upload.MimeType)
	if upload.FileType == TypeUnknown {
		return nil, fmt.Errorf("Unkown file type:%v \n", upload.MimeType)
	}
	err = validateMimeType(upload.FileType, upload.MimeType)
	if err != nil {
		return nil, err
	}
	upload.Extension = path.Ext(file.Filename)
	err = validateExtensions(upload.FileType, upload.Extension)
	if err != nil {
		return nil, err
	}
	upload.Size = file.Size
	err = validateMaxSize(upload.FileType, upload.Size)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	date := now.Format("2006-01-02")
	guid := random.Uuid(false)
	if subDir == SubDir {
		subDir = upload.FileType
	}
	upload.File = subDir + "/" + date + "/" + guid + upload.Extension
	filePath := GetPath() + upload.File
	filePath = filepath.ToSlash(filePath)
	err = mkdir(path.Dir(filePath))
	if err != nil {
		return nil, err
	}
	err = ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		return nil, err
	}
	return upload, nil
}

// Base64Upload 图片base64上传
func Base64Upload(base64Content string, subDir string) (*Upload, error) {
	upload := &Upload{}
	//验证图片格式
	matched, _ := regexp.MatchString(`^data:\s*image\/(\w+);base64,`, base64Content)
	if !matched {
		return nil, errors.New("数据格式错误")
	}
	//解析图片mimeType、extension和内容
	reg, _ := regexp.Compile(`^data:(image\/(\w+));base64,`)
	allData := reg.FindAllSubmatch([]byte(base64Content), 2)
	upload.MimeType = string(allData[0][1]) //image/png
	upload.FileType = "image"
	upload.Extension = "." + string(allData[0][2]) //png ，jpeg 后缀获取
	base64Str := reg.ReplaceAllString(base64Content, "")
	//拼接图片保存路径
	now := time.Now()
	date := now.Format("2006-01-02")
	guid := random.Uuid(false)
	if subDir == SubDir {
		subDir = TypeImage
	}
	upload.File = subDir + "/" + date + "/" + guid + upload.Extension
	filePath := GetPath() + upload.File
	filePath = filepath.ToSlash(filePath)
	//创建文件夹
	err := mkdir(path.Dir(filePath))
	if err != nil {
		return nil, err
	}
	//写入文件
	err = ioutil.WriteFile(filePath, base64.Decode(base64Str), os.ModePerm)
	if err != nil {
		return nil, err
	}
	//获取文件大小
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	upload.Size = stat.Size()
	return upload, nil
}

func GetPath() string {
	return config.Config.Web.Upload.Path
}

func GetUrl() string {
	return config.Config.Web.Upload.Url
}

// GetFileType 获取文件类型
func GetFileType(mimeType string) string {
	matched, _ := regexp.MatchString(`^image\/(\w+)`, mimeType)
	if matched {
		return TypeImage
	}
	matched, _ = regexp.MatchString(`^video\/(\w+)`, mimeType)
	if matched {
		return TypeVideo
	}
	return TypeUnknown
}

func Remove(fullPath string) error {
	path := UrlToPath(Url(fullPath))
	if path == "" {
		return nil
	}
	filePath := string(path)
	_, err := os.Stat(filePath)
	if err != nil {
		log.Debugf(err.Error())
		return nil
	}
	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func FullUrl(relativeUrl string) string {
	if relativeUrl == "" {
		return ""
	}
	urlPrefix := GetUrl()
	if strings.HasPrefix(relativeUrl, urlPrefix) {
		return relativeUrl
	}
	url := PathToUrl(Path(relativeUrl))
	return urlPrefix + string(url)
}

func FullPath(relativePath string) string {
	if relativePath == "" {
		return ""
	}
	pathPrefix := GetPath()
	if strings.HasPrefix(relativePath, pathPrefix) {
		return relativePath
	}
	path := UrlToPath(Url(relativePath))
	return pathPrefix + string(path)
}

func PathToUrl(path Path) Url {
	if path == "" {
		return ""
	}
	pathString := string(path)
	pathString = filepath.FromSlash(pathString)
	pathString = strings.Replace(pathString, GetPath(), GetUrl(), 1)
	return Url(filepath.ToSlash(pathString))
}

func UrlToPath(url Url) Path {
	if url == "" {
		return ""
	}
	urlString := string(url)
	urlString = filepath.ToSlash(urlString)
	urlString = strings.Replace(urlString, GetUrl(), GetPath(), 1)
	return Path(filepath.FromSlash(urlString))
}

func RelativeUrl(url Url) string {
	return strings.TrimPrefix(string(url), GetUrl())
}

func RelativePath(path Path) string {
	return strings.TrimPrefix(string(path), GetPath())
}

func validateExtensions(fileType string, fileExtension string) error {
	extensions := make([]string, 0)
	switch fileType {
	case TypeVideo:
		extensions = config.Config.Web.Upload.Video.Extensions
	case TypeImage:
		extensions = config.Config.Web.Upload.Image.Extensions
	}
	for _, extension := range extensions {
		if extension == "*" {
			return nil
		}
		if extension == fileExtension {
			return nil
		}
		if extension == ("." + fileExtension) {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Not supperted extension: %v/%v \n", fileExtension, extensions))
}

func validateMimeType(fileType string, mimeType string) error {
	mimeTypes := make([]string, 0)
	switch fileType {
	case TypeVideo:
		mimeTypes = config.Config.Web.Upload.Video.MimeTypes
	case TypeImage:
		mimeTypes = config.Config.Web.Upload.Image.MimeTypes
	}
	for _, mimeType := range mimeTypes {
		if mimeType == "*" {
			return nil
		}
		if mimeType == mimeType {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Not supperted mime type: %v/%v \n", mimeType, mimeType))
}

func validateMaxFiles(fileType string, fileCount int) error {
	maxFiles := 1
	switch fileType {
	case TypeVideo:
		maxFiles = config.Config.Web.Upload.Video.MaxFiles
	case TypeImage:
		maxFiles = config.Config.Web.Upload.Image.MaxFiles
	}
	if maxFiles >= fileCount {
		return nil
	}
	return errors.New(fmt.Sprintf("Files count is exceeds limit: %v/%v \n", fileCount, maxFiles))
}

func validateMaxSize(fileType string, fileSize int64) error {
	maxSize := ""
	switch fileType {
	case TypeVideo:
		maxSize = config.Config.Web.Upload.Video.MaxSize
	case TypeImage:
		maxSize = config.Config.Web.Upload.Image.MaxSize
	}
	if maxSize == "" {
		return nil
	}
	maxByte := sizeToBytes(maxSize)
	if maxByte >= fileSize {
		return nil
	}
	return errors.New(fmt.Sprintf("File size is exceeds limit: %v/%v \n", fileSize, maxSize))
}

func sizeToBytes(sizeSetting interface{}) int64 {
	switch sizeSetting.(type) {
	case int:
		return int64(sizeSetting.(int))
	case string:
		var strategy int
		sizeStr := sizeSetting.(string)
		strLen := len(sizeStr)
		lastChar := sizeStr[strLen-1:]
		switch lastChar {
		case "M", "m":
			strategy = 1048576
		case "K", "k":
			strategy = 1024
		case "G", "g":
			strategy = 1073741824
		default:
			strategy = 1
		}
		sizeStr = sizeStr[:strLen-1]
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil {
			log.Errorf("Unable to covert max size:%v %e", sizeSetting, err)
			return 0
		}
		return int64(sizeInt * strategy)
	}
	return 0
}

func mkdir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		return nil
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err == nil {
		return nil
	}
	return err
}
