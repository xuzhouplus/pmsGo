package config

type session struct {
	Name   string `yaml:"name"`
	Secret string `yaml:"secret"`
	Type   string `yaml:"type"`
	Idle   int `yaml:"idle"`
}
