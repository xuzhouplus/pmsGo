package file

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"math"
	"mime"
	"os"
	"path/filepath"
	"strconv"
)

type Image struct {
	File      image.Image
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mimeType"`
	Extension string `json:"extension"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func New(file string) (*Image, error) {
	open, err := imaging.Open(file)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	return &Image{File: open, Path: filepath.Dir(file), Name: stat.Name(), Extension: filepath.Dir(file), Size: stat.Size(), MimeType: mime.TypeByExtension(filepath.Ext(file)), Width: open.Bounds().Size().X, Height: open.Bounds().Size().Y}, nil
}
func (img Image) Resize(width int, height int) (*Image, error) {
	target := img.Path + img.Name + strconv.Itoa(width) + "_" + strconv.Itoa(height) + ".jpg"
	resize := imaging.Resize(img.File, width, height, imaging.Lanczos)
	dst := imaging.New(img.Width, img.Height, color.NRGBA{})
	dst = imaging.Paste(dst, resize, image.Pt(0, 0))
	err := imaging.Save(dst, target)
	if err != nil {
		return nil, err
	}
	rst, err := New(target)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (img Image) Cut(offsetX int, offsetY int, width int, height int) (*Image, error) {
	target := img.Path + img.Name + strconv.Itoa(width) + "_" + strconv.Itoa(height) + ".jpg"
	crop := imaging.Crop(img.File, image.Rectangle{Min: image.Pt(offsetX, offsetY), Max: image.Pt(int(math.Min(float64(offsetX+width), float64(img.Width))), int(math.Min(float64(offsetY+height), float64(img.Height))))})
	dst := imaging.New(width, height, color.NRGBA{})
	dst = imaging.Paste(dst, crop, image.Pt(0, 0))
	err := imaging.Save(dst, target)
	if err != nil {
		return nil, err
	}
	rst, err := New(target)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (img Image) Compress(quality int) {

}
func (img Image) Convert(extension string) {

}

func (img Image) Blur(sigma float64) (*Image, error) {
	target := img.Path + img.Name + ".blur.jpg"
	blur := imaging.Blur(img.File, sigma)
	dst := imaging.New(img.Width, img.Height, color.NRGBA{})
	dst = imaging.Paste(dst, blur, image.Pt(0, 0))
	err := imaging.Save(blur, target)
	if err != nil {
		return nil, err
	}
	rst, err := New(target)
	if err != nil {
		return nil, err
	}
	return rst, nil
}
