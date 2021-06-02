package config

type log struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
	Json  bool   `yaml:"json"`
}
