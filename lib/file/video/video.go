package video

import (
	"github.com/floostack/transcoder"
	"github.com/floostack/transcoder/ffmpeg"
	"mime"
	"os"
	"path/filepath"
	"pmsGo/lib/config"
	"pmsGo/lib/file"
	"pmsGo/lib/log"
	"strconv"
	"strings"
)

// https://www.cnblogs.com/2020-zhy-jzoj/p/13165086.html

type Video struct {
	file.File
	ffmpeg   transcoder.Transcoder
	BitRate  string `json:"bitRate""`
	Duration string `json:"duration"`
}

const (
	MimeTypeMp4 = "video/mp4"
)

func Open(file string) (*Video, error) {
	ffmpegConf := &ffmpeg.Config{
		FfmpegBinPath:   config.Config.Web.Upload.Video.Ffmpeg,
		FfprobeBinPath:  config.Config.Web.Upload.Video.Ffprobe,
		ProgressEnabled: true,
	}
	video := &Video{}
	video.ffmpeg = ffmpeg.New(ffmpegConf).Input(file)
	video.Name = filepath.Base(file)
	video.Path = filepath.Dir(file)
	video.Extension = filepath.Ext(file)
	video.MimeType = mime.TypeByExtension(filepath.Ext(file))
	err := video.Metadata()
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (receiver *Video) Metadata() error {
	metadata, err := receiver.ffmpeg.GetMetadata()
	if err == nil {
		format := metadata.GetFormat()
		receiver.BitRate = format.GetBitRate()
		receiver.Duration = format.GetDuration()
		receiver.Size = format.GetSize()
		log.Debugf("%+v", format)
		for _, stream := range metadata.GetStreams() {
			if stream.GetCodecType() == "video" {
				log.Debugf("%+v", stream)
				receiver.Height = stream.GetHeight()
				receiver.Width = stream.GetWidth()
				break
			}
		}
		return nil
	} else {
		return err
	}
}
func (receiver Video) FileName() string {
	return strings.TrimSuffix(receiver.Name, receiver.Extension)
}

func (receiver Video) CreateThumb(width int, height int, ext string, time string) (string, error) {
	if time == "" {
		time = "00:00:00"
	}
	widthSpread := strconv.Itoa(width)
	heightSpread := strconv.Itoa(height)
	path := receiver.Path + string(filepath.Separator) + receiver.FileName() + string(filepath.Separator) + "_" + widthSpread + "_" + heightSpread + "." + ext
	dir := filepath.Dir(path)
	os.Mkdir(dir, 755)
	videoProfile := "baseline"
	outputFormat := "image2"
	overwrite := true
	videoFilter := "scale=" + widthSpread + ":" + heightSpread
	vFrames := 1
	opts := ffmpeg.Options{
		VideoProfile: &videoProfile,
		OutputFormat: &outputFormat,
		SeekTime:     &time,
		Overwrite:    &overwrite,
		Vframes:      &vFrames,
		VideoFilter:  &videoFilter,
	}
	_, err := receiver.ffmpeg.
		Output(path).
		WithOptions(opts).
		Start(opts)
	if err != nil {
		return "", err
	}
	return path, nil
}

type Progress func(<-chan transcoder.Progress)

func (receiver Video) CreateM3u8(width int, height int, progress Progress) (string, error) {
	widthSpread := strconv.Itoa(width)
	heightSpread := strconv.Itoa(height)
	path := receiver.Path + string(filepath.Separator) + receiver.FileName() + string(filepath.Separator) + "_" + widthSpread + "_" + heightSpread + ".m3u8"
	dir := filepath.Dir(path)
	os.Mkdir(dir, 755)
	videoProfile := "baseline"
	hlsListSize := 0
	hlsSegmentFilename := dir + string(filepath.Separator) + widthSpread + "_" + heightSpread + "%05d.ts"
	hlsSegmentDuration := 10
	outputFormat := "hls"
	overwrite := true
	videoFilter := "scale=" + widthSpread + ":" + heightSpread
	opts := ffmpeg.Options{
		VideoProfile:       &videoProfile,
		HlsListSize:        &hlsListSize,
		HlsSegmentFilename: &hlsSegmentFilename,
		HlsSegmentDuration: &hlsSegmentDuration,
		OutputFormat:       &outputFormat,
		Overwrite:          &overwrite,
		VideoFilter:        &videoFilter,
		ExtraArgs: map[string]interface{}{
			"-level":        "3.0",
			"-start_number": "0",
		},
	}
	progressChannel, err := receiver.ffmpeg.
		Output(path).
		WithOptions(opts).
		Start(opts)
	if err != nil {
		return "", err
	}
	progress(progressChannel)
	return path, nil
}