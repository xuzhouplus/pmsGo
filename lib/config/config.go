package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	Site     site
	Database database
	Redis    redis
	Cache    cache
	Session  session
	Log      log
	Sync     sync
	Web      Web
}

var Config = &config{}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dirSep := string(filepath.Separator)
	cfgFile := pwd + dirSep + "config" + dirSep + "app.yaml"
	file, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(fmt.Sprintf("Can`t read config image:%err\n", err))
	}

	err = yaml.Unmarshal(file, Config)
	if err != nil {
		panic(fmt.Sprintf("Can`t analyse config image:%err\n", err))
	}
}
