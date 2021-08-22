package config

type site struct {
	Name     string `yaml:"name"`
	Debug    bool   `yaml:"debug"`
	Maintain bool   `yaml:"maintain"`
	Listen     string    `yaml:"listen"`
	AllowCrossDomain bool `yaml:"cors"`
}
