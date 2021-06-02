package config

type cache struct {
	Prefix string `yaml:"prefix"`
	Expire int    `yaml:"expire"`
}
