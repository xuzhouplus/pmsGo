package image

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/png"
	"math"
	"mime"
	"os"
	"path/filepath"
	"pmsGo/lib/file"
	"strconv"
	"strings"
)

const (
	MimeTypeJpg  = "image/jpeg"
	MimeTypeJpeg = "image/jpeg"
	MimeTypePng  = "image/png"
)

type Image struct {
	Image image.Image
	file.File
}

func Open(file string) (*Image, error) {
	open, err := imaging.Open(file, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	img := &Image{}
	img.Image = open
	img.Path = filepath.Dir(file)
	img.Name = stat.Name()
	img.Extension = filepath.Ext(file)
	img.Size = string(stat.Size())
	img.MimeType = mime.TypeByExtension(filepath.Ext(file))
	img.Width = open.Bounds().Size().X
	img.Height = open.Bounds().Size().Y
	return img, nil
}

func (img Image) FileName() string {
	return strings.TrimSuffix(img.Name, img.Extension)
}

func (img Image) FullPath() string {
	return img.Path + string(filepath.Separator) + img.Name
}

func (img Image) Resize(width int, height int) *image.NRGBA {
	return imaging.Resize(img.Image, width, height, imaging.Lanczos)
}

func (img Image) Cut(offsetX int, offsetY int, width int, height int) *image.NRGBA {
	return imaging.Crop(img.Image, image.Rectangle{Min: image.Pt(offsetX, offsetY), Max: image.Pt(int(math.Min(float64(offsetX+width), float64(img.Width))), int(math.Min(float64(offsetY+height), float64(img.Height))))})
}

func (img Image) CompressJPEG(quality int) (imaging.EncodeOption, error) {
	if img.MimeType != MimeTypeJpeg {
		return nil, fmt.Errorf("not supported image format：%v", img.MimeType)
	}
	return imaging.JPEGQuality(quality), nil
}

func (img Image) CompressPNG(level png.CompressionLevel) (imaging.EncodeOption, error) {
	if img.MimeType != MimeTypePng {
		return nil, fmt.Errorf("not supported image format：%v", img.MimeType)
	}
	return imaging.PNGCompressionLevel(level), nil
}

func (img Image) Convert(extension string) (*Image, error) {
	target := img.Path + img.Name + extension
	err := imaging.Save(img.Image, target)
	if err != nil {
		return nil, err
	}
	rst, err := Open(target)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (img Image) Blur(sigma float64) *image.NRGBA {
	return imaging.Blur(img.Image, sigma)
}

func (img Image) Save(filePath string, opt imaging.EncodeOption) (*Image, error) {
	err := imaging.Save(img.Image, filePath, opt)
	if err != nil {
		return nil, err
	}
	rst, err := Open(filePath)
	if err != nil {
		return nil, err
	}
	return rst, nil
}

func (img Image) CreateCarousel(destName string, width int, height int, extension string) (*Image, error) {
	if extension == "" {
		extension = strings.TrimPrefix(img.Extension, ".")
	}
	blur := img.Blur(3)
	bg := imaging.Resize(blur, width, height, imaging.Lanczos)
	fg := img.Image
	if img.Width > width || img.Height > height {
		fg = imaging.Fit(img.Image, width, height, imaging.Lanczos)
	}
	dst := imaging.OverlayCenter(bg, fg, 1.0)
	tg := img.Path + string(filepath.Separator) + destName + "." + extension
	err := imaging.Save(dst, tg)
	if err != nil {
		return nil, err
	}
	carousel, err := Open(tg)
	if err != nil {
		return nil, err
	}
	return carousel, nil
}

func (img Image) CreateResize(destName string, width int, height int, ext string) (*Image, error) {
	if ext == "" {
		ext = strings.TrimPrefix(img.Extension, ".")
	}
	fg := imaging.Fit(img.Image, width, height, imaging.Lanczos)
	var bg *image.NRGBA
	if ext == "png" {
		bg = imaging.New(width, height, color.NRGBA{})
	} else {
		bg = imaging.New(width, height, color.White)
	}
	dst := imaging.OverlayCenter(bg, fg, 1)
	path := img.Path + string(filepath.Separator) + destName + "_" + strconv.Itoa(width) + "_" + strconv.Itoa(height) + "." + ext
	err := imaging.Save(dst, path)
	if err != nil {
		return nil, err
	}
	target, err := Open(path)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func (img Image) CreateCompress(destName string, quality int) (*Image, error) {
	switch img.MimeType {
	case MimeTypePng:
		target := img.Path + string(filepath.Separator) + destName + "_" + strconv.Itoa(quality) + img.Extension
		level := png.DefaultCompression
		if quality >= 75 && quality < 100 {
			level = png.NoCompression
		} else if quality >= 50 && quality < 75 {
			level = png.DefaultCompression
		} else if quality >= 25 && quality < 50 {
			level = png.BestSpeed
		} else {
			level = png.BestCompression
		}
		encodeOptions, err := img.CompressPNG(level)
		if err != nil {
			return nil, err
		}
		save, err := img.Save(target, encodeOptions)
		if err != nil {
			return nil, err
		}
		return save, nil
	case MimeTypeJpg:
		target := img.Path + string(filepath.Separator) + img.FileName() + "_" + strconv.Itoa(quality) + img.Extension
		encodeOptions, err := img.CompressJPEG(quality)
		if err != nil {
			return nil, err
		}
		save, err := img.Save(target, encodeOptions)
		if err != nil {
			return nil, err
		}
		return save, nil
	default:
		return nil, fmt.Errorf("unsupported image format: %v", img.MimeType)
	}
}
