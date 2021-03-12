package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type config struct {
	Site     site
	Database database
	Redis    redis
}

var Config = &config{}

func init() {
	pwd, error := os.Getwd()
	if error != nil {
		fmt.Println(error)
		return
	}
	file, error := ioutil.ReadFile(pwd + "\\config\\app.yaml")
	if error != nil {
		fmt.Printf("Can`t read config file:%err\n", error)
		return
	}

	error = yaml.Unmarshal(file, Config)
	if error != nil {
		fmt.Errorf("Can`t analyse config file:%err\n", error)
	}
}
