package config

type Auth struct {
	Except   []string `yaml:"except"`
	Optional []string `yaml:"optional"`
}
type Upload struct {
	Path       string   `yaml:"path"`
	Url        string   `yaml:"url"`
	Extensions []string `yaml:"extensions"`
	MaxSize    []string `yaml:"maxSize"`
	MaxFiles   int      `yaml:"maxFiles"`
	MimeTypes  []string `yaml:"mimeTypes"`
}
type Web struct {
	Host     string            `yaml:"host"`
	Connects []string          `yaml:"connects"`
	Security map[string]string `yaml:"security"`
	Auth     map[string]Auth   `yaml:"auth"`
	Upload   Upload            `yaml:"upload"`
}
