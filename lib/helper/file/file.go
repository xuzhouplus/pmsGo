package file

import (
	"mime"
	"os"
	"path/filepath"
)

type File struct {
	File      *os.File
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mimeType"`
	Extension string `json:"extension"`
}

func New(file string) (*File, error) {
	open, err := os.Open(file)
	if err != nil {
		return &File{}, err
	}
	defer open.Close()
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	return &File{File: open, Path: filepath.Dir(file), Name: filepath.Base(file), Extension: filepath.Dir(file), Size: stat.Size(), MimeType: mime.TypeByExtension(filepath.Ext(file))}, nil
}
func (file File) Resize() {

}
