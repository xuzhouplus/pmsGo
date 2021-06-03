package config

type log struct {
	File  string `yaml:"image"`
	Level string `yaml:"level"`
	Json  bool   `yaml:"json"`
}
