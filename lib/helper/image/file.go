package image

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pmsGo/lib/config"
	"strconv"
	"strings"
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

type Url string

type Path string

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

func Remove(fullPath string) error {
	path := UrlToPath(Url(fullPath))
	filePath := string(path)
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func FullUrl(relativeUrl string) string {
	if strings.HasPrefix(relativeUrl, Settings.url) {
		return relativeUrl
	}
	url := PathToUrl(Path(relativeUrl))
	return Settings.url + string(url)
}

func FullPath(relativePath string) string {
	if strings.HasPrefix(relativePath, Settings.path) {
		return relativePath
	}
	path := UrlToPath(Url(relativePath))
	return Settings.path + string(path)
}

func PathToUrl(path Path) Url {
	pathString := string(path)
	pathString = filepath.FromSlash(pathString)
	pathString = strings.Replace(pathString, Settings.path, Settings.url, 1)
	return Url(filepath.ToSlash(pathString))
}

func UrlToPath(url Url) Path {
	urlString := string(url)
	urlString = filepath.ToSlash(urlString)
	urlString = strings.Replace(urlString, Settings.url, Settings.path, 1)
	return Path(filepath.FromSlash(urlString))
}

func RelativeUrl(url Url) string {
	return strings.TrimPrefix(string(url), Settings.url)
}

func RelativePath(path Path) string {
	return strings.TrimPrefix(string(path), Settings.path)
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
