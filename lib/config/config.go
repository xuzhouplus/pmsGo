package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	Site     site
	Database database
	Redis    redis
	Session  session
	Web      Web
}

var Config = &config{}

func init() {
	pwd, error := os.Getwd()
	if error != nil {
		panic(error)
	}
	dirSep := string(filepath.Separator)
	cfgFile := pwd + dirSep + "config" + dirSep + "app.yaml"
	log.Println(fmt.Sprintf("Load config file:%v", cfgFile))
	file, error := ioutil.ReadFile(cfgFile)
	if error != nil {
		panic(fmt.Sprintf("Can`t read config file:%err\n", error))
	}

	error = yaml.Unmarshal(file, Config)
	if error != nil {
		panic(fmt.Sprintf("Can`t analyse config file:%err\n", error))
	}
}
