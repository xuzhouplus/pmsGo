package file

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
)

type File struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      string `json:"size"`
	MimeType  string `json:"mimeType"`
	Extension string `json:"extension"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type Upload struct {
	ctx       *gin.Context
	Uuid      string `json:"uuid"`
	File      string `json:"file"`
	FileType  string `json:"fileType"`
	MimeType  string `json:"mimeType"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

// Path 获取文件绝对路径
func (upload Upload) Path() Path {
	return Path(filepath.FromSlash(GetPath() + upload.File))
}

// Url 获取文件访问Url
func (upload Upload) Url() Url {
	return Url(filepath.ToSlash(GetUrl() + upload.File))
}
