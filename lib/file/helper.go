package file

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"pmsGo/lib/cache"
	"pmsGo/lib/config"
	"pmsGo/lib/log"
	"pmsGo/lib/security/base64"
	"pmsGo/lib/security/random"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	TypeVideo   = "video"
	TypeImage   = "image"
	TypeUnknown = "unknown"
)

const (
	B  = 1
	KB = B * 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
)

const SubDir = "baseOnFileType"

const ChunkUploadCountCache = "chunk:count:"
const ChunkUploadFileCache = "chunk:file:"

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
	upload.Extension = filepath.Ext(file.Filename)
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
	upload.Uuid = guid
	if subDir == SubDir {
		subDir = upload.FileType
	}
	upload.File = "/" + subDir + "/" + date + "/" + guid + "/" + file.Filename
	filePath := GetPath() + upload.File
	filePath = filepath.ToSlash(filePath)
	err = Mkdir(filepath.Dir(filePath))
	if err != nil {
		return nil, err
	}
	err = ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		return nil, err
	}
	return upload, nil
}

func ChunkCheck(ctx *gin.Context) bool {
	taskId := ctx.Query("id")
	chunkIndex := ctx.Query("index")
	cached, err := cache.Get(ChunkUploadCountCache + taskId + ":" + chunkIndex)
	if err != nil {
		return false
	}
	if cached == "1" {
		return true
	}
	return false
}

// ChunkUpload Resumablejs分片上传
func ChunkUpload(ctx *gin.Context, fieldName string, subDir string) (*Upload, error) {
	upload := &Upload{}
	upload.ctx = ctx
	upload.MimeType = ctx.PostForm("type")
	upload.FileType = GetFileType(upload.MimeType)
	fileName := ctx.PostForm("file")
	upload.Extension = filepath.Ext(fileName)
	upload.Uuid = ctx.PostForm("id")
	if subDir == SubDir {
		subDir = upload.FileType
	}
	lock := &sync.Mutex{}
	lock.Lock()
	cacheData, err := cache.Get(ChunkUploadFileCache + upload.Uuid)
	if err != nil {
		return nil, err
	}
	var saveFile *os.File
	var filePath string
	if cacheData == nil {
		now := time.Now()
		date := now.Format("2006-01-02")
		upload.File = "/" + subDir + "/" + date + "/" + upload.Uuid + "/source" + upload.Extension
		filePath = GetPath() + upload.File
		filePath = filepath.FromSlash(filePath)
		total, err := strconv.Atoi(ctx.PostForm("total"))
		if err != nil {
			return nil, err
		}
		saveFile, err = CreateFile(filePath, int64(total))
		if err != nil {
			return nil, err
		}
		defer saveFile.Close()
		err = cache.Set(ChunkUploadFileCache+upload.Uuid, filePath, 1800)
		if err != nil {
			return nil, err
		}
	} else {
		filePath = cacheData.(string)
		saveFile, err = OpenFile(filePath)
		if err != nil {
			return nil, err
		}
		defer saveFile.Close()
		upload.File = RelativePath(Path(filePath))
	}
	lock.Unlock()
	formFile, err := ctx.FormFile(fieldName)
	if err != nil {
		return nil, err
	}
	uploadFile, err := formFile.Open()
	if err != nil {
		return nil, err
	}
	defer uploadFile.Close()
	chunkIndex, _ := strconv.Atoi(ctx.PostForm("index"))
	chunkSize, _ := strconv.Atoi(ctx.PostForm("chunk"))
	chunkOffset := (chunkIndex - 1) * chunkSize * B
	bufferSize := make([]byte, 1024)
	for {
		readLen, err := uploadFile.Read(bufferSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		_, err = saveFile.WriteAt(bufferSize, int64(chunkOffset))
		if err != nil {
			return nil, err
		}
		chunkOffset = chunkOffset + readLen
	}
	increase, err := cache.Increase(ChunkUploadCountCache+upload.Uuid, 1)
	if err != nil {
		return nil, err
	}
	lock.Lock()
	defer lock.Unlock()
	if strconv.Itoa(int(increase)) == ctx.PostForm("count") {
		upload.Status = UploadStatusComplete
	} else {
		upload.Status = UploadStatusProcess
	}
	return upload, nil
}

// Base64Upload 图片base64上传
func Base64Upload(base64Content string, subDir string) (*Upload, error) {
	upload := &Upload{}
	//验证图片格式
	matched, _ := regexp.MatchString(`^data:\s*image/(\w+);base64,`, base64Content)
	if !matched {
		return nil, errors.New("数据格式错误")
	}
	//解析图片mimeType、extension和内容
	reg, _ := regexp.Compile(`^data:(image/(\w+));base64,`)
	allData := reg.FindAllSubmatch([]byte(base64Content), 2)
	upload.MimeType = string(allData[0][1]) //image/png
	upload.FileType = "image"
	upload.Extension = "." + string(allData[0][2]) //png ，jpeg 后缀获取
	base64Str := reg.ReplaceAllString(base64Content, "")
	//拼接图片保存路径
	now := time.Now()
	date := now.Format("2006-01-02")
	guid := random.Uuid(false)
	upload.Uuid = guid
	if subDir == SubDir {
		subDir = TypeImage
	}
	upload.File = subDir + "/" + date + "/" + guid + "/" + "source" + upload.Extension
	filePath := GetPath() + upload.File
	filePath = filepath.ToSlash(filePath)
	//创建文件夹
	err := Mkdir(filepath.Dir(filePath))
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
	return filepath.FromSlash(config.Config.Web.Upload.Path)
}

func GetUrl() string {
	return filepath.ToSlash(config.Config.Web.Upload.Url)
}

// GetFileType 获取文件类型
func GetFileType(mimeType string) string {
	matched, _ := regexp.MatchString(`^image/(\w+)`, mimeType)
	if matched {
		return TypeImage
	}
	matched, _ = regexp.MatchString(`^video/(\w+)`, mimeType)
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
	pathString = filepath.ToSlash(pathString)
	pathString = strings.Replace(pathString, GetPath(), GetUrl(), 1)
	return Url(filepath.ToSlash(pathString))
}

func UrlToPath(url Url) Path {
	if url == "" {
		return ""
	}
	urlString := string(url)
	urlString = filepath.FromSlash(urlString)
	urlString = strings.Replace(urlString, GetUrl(), GetPath(), 1)
	return Path(filepath.FromSlash(urlString))
}

func RelativeUrl(url Url) string {
	return strings.TrimPrefix(string(url), GetUrl())
}

func RelativePath(path Path) string {
	log.Debug(path)
	log.Debug(GetPath())
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
		if ("." + extension) == fileExtension {
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
	for _, rangeMimeType := range mimeTypes {
		if rangeMimeType == "*" {
			return nil
		}
		if rangeMimeType == mimeType {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Not supperted mime type: %v / %v \n", mimeType, mimeTypes))
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

func Mkdir(path string) error {
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

func CreateFile(file string, size int64) (*os.File, error) {
	err := Mkdir(filepath.Dir(file))
	if err != nil {
		return nil, err
	}
	openFile, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	if err := openFile.Truncate(size * B); err != nil {
		err := openFile.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return openFile, nil
}

func OpenFile(file string) (*os.File, error) {
	return os.OpenFile(file, os.O_RDWR, 0777)
}
