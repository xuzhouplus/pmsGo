package config

type Video struct {
	Ffmpeg     string   `yaml:"ffmpeg"`
	Ffprobe    string   `yaml:"ffprobe"`
	Extensions []string `yaml:"extensions"`
	MaxSize    string   `yaml:"maxSize"`
	MaxFiles   int      `yaml:"maxFiles"`
	MimeTypes  []string `yaml:"mimeTypes"`
}
type Image struct {
	Extensions []string `yaml:"extensions"`
	MaxSize    string   `yaml:"maxSize"`
	MaxFiles   int      `yaml:"maxFiles"`
	MimeTypes  []string `yaml:"mimeTypes"`
}
type Upload struct {
	Path  string `yaml:"path"`
	Url   string `yaml:"url"`
	Video Video
	Image Image
}
type Web struct {
	Host     string            `yaml:"host"`
	Connects []string          `yaml:"connects"`
	Security map[string]string `yaml:"security"`
	Upload   Upload            `yaml:"upload"`
}
