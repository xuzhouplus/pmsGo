package file

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
	"path/filepath"
	"pmsGo/lib/config"
	"pmsGo/lib/security"
	"strconv"
	"strings"
	"time"
)

type settings struct {
	path       string
	url        string
	extensions []string
	maxSize    int64
	maxFiles   int
	mimeTypes  []string
}

var Settings = &settings{}

func init() {
	settings := config.Config.Web["upload"].(map[interface{}]interface{})
	Settings.path = filepath.FromSlash(settings["path"].(string))
	Settings.url = filepath.ToSlash(settings["url"].(string))
	extensions := settings["extensions"].([]interface{})
	for _, extension := range extensions {
		Settings.extensions = append(Settings.extensions, extension.(string))
	}
	Settings.maxSize = sizeToBytes(settings["maxSize"])
	Settings.maxFiles = settings["maxFiles"].(int)
	mimeTypes := settings["mimeTypes"].([]interface{})
	for _, mimeType := range mimeTypes {
		Settings.mimeTypes = append(Settings.mimeTypes, mimeType.(string))
	}
	log.Printf("Upload: %v \n", Settings)
}

type Upload struct {
	ctx       *gin.Context
	File      string `json:"file"`
	MimeType  string `json:"mimeType"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

func (helper *Upload) Upload(ctx *gin.Context, fieldName string, subDir string) error {
	helper.ctx = ctx
	file, err := ctx.FormFile(fieldName)
	if err != nil {
		return err
	}
	log.Printf("Upload file:%v \n", file.Header)
	helper.MimeType = file.Header.Get("Content-Type")
	err = validateMimeType(helper.MimeType)
	if err != nil {
		return err
	}
	helper.Extension = path.Ext(file.Filename)
	err = validateExtensions(helper.Extension)
	if err != nil {
		return err
	}
	helper.Size = file.Size
	err = validateMaxSize(helper.Size)
	if err != nil {
		return err
	}
	now := time.Now()
	date := now.Format("2006-01-02")
	guid := security.Uuid(false)
	helper.File = subDir + "/" + date + "/" + guid + path.Ext(file.Filename)
	filePath := Settings.path + helper.File
	filePath = filepath.ToSlash(filePath)
	err = mkdir(path.Dir(filePath))
	if err != nil {
		return err
	}
	err = ctx.SaveUploadedFile(file, filePath)
	if err != nil {
		return err
	}
	return nil
}

func (helper Upload) Path() string {
	return filepath.FromSlash(Settings.path + helper.File)
}

func (helper Upload) Url() string {
	return filepath.ToSlash(Settings.url + helper.File)
}

func PathToUrl(path string) string {
	path = filepath.FromSlash(path)
	path = strings.Replace(path, Settings.path, Settings.url, 1)
	return filepath.ToSlash(path)
}

func UrlToPath(url string) string {
	url = filepath.ToSlash(url)
	url = strings.Replace(url, Settings.url, Settings.path, 1)
	return filepath.FromSlash(url)
}

func validateExtensions(fileExtension string) error {
	for _, extension := range Settings.extensions {
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
	return errors.New(fmt.Sprintf("Not supperted extension: %v/%v \n", fileExtension, Settings.extensions))
}

func validateMimeType(mimeType string) error {
	for _, mimeType := range Settings.mimeTypes {
		if mimeType == "*" {
			return nil
		}
		if mimeType == mimeType {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Not supperted mime type: %v/%v \n", mimeType, Settings.mimeTypes))
}

func validateMaxFiles(fileCount int) error {
	if Settings.maxFiles >= fileCount {
		return nil
	}
	return errors.New(fmt.Sprintf("Files count is exceeds limit: %v/%v \n", fileCount, Settings.maxFiles))
}

func validateMaxSize(fileSize int64) error {
	if Settings.maxSize == 0 {
		return nil
	}
	if Settings.maxSize >= fileSize {
		return nil
	}
	return errors.New(fmt.Sprintf("File size is exceeds limit: %v/%v \n", fileSize, Settings.maxSize))
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
			log.Printf("Unable to covert max size:%v %e", sizeSetting, err)
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
