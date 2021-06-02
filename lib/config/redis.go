package config

type redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Auth     string `yaml:"auth"`
	Database int    `yaml:"database"`
}
