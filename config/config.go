package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sync"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

var instance *Config
var once sync.Once

func Get() *Config {
	return instance
}

func Init(e string) {
	once.Do(func() {
		env := e
		if e == "" {
			env = "production"
		}

		p, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		var filePath string
		filePath = p + "/config/database.yml"

		var conf map[string]Config
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(buf, &conf)
		if err != nil {
			panic(err)
		}

		instance = &Config{
			Username: conf[env].Username,
			Password: conf[env].Password,
			Host:     conf[env].Host,
			Port:     conf[env].Port,
			Database: conf[env].Database,
		}
	})
}
