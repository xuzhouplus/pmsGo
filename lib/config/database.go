package config

type database struct {
	Host            string `yaml:"host"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	Prefix          string `yaml:"prefix"`
	MaxIdleConnect  int    `yaml:"max_idle_connect"`
	MaxOpenConnect  int    `yaml:"max_open_connect"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}
