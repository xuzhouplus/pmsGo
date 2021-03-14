package config

type site struct {
	Name     string `yaml:"name"`
	Debug    bool   `yaml:"debug"`
	Maintain bool   `yaml:"maintain"`
	Port     int    `yaml:"port"`
	AllowCrossDomain bool `yaml:"cors"`
}
