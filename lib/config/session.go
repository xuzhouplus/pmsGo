package config

type session struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
	Secret string `yaml:"secret"`
	Type   string `yaml:"type"`
	Idle   int    `yaml:"idle"`
}
