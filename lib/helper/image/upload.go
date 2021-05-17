package image

import (
	"github.com/gin-gonic/gin"
	"log"
	"path"
	"path/filepath"
	"pmsGo/lib/security"
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
	helper.File = subDir + "/" + date + "/" + guid + path.Ext(file.Filename)
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

func (helper Instance) Path() Path {
	return Path(filepath.FromSlash(Settings.path + helper.File))
}

func (helper Instance) Url() Url {
	return Url(filepath.ToSlash(Settings.url + helper.File))
}
