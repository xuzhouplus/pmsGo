package config

type database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset string `yaml:"charset"`
	Prefix   string `yaml:"prefix"`
}