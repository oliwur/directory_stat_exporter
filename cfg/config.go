package cfg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Dir struct {
	Path      string
	Recursive bool
}
type Config struct {
	Directories []Dir
}

func GetConfig() Config {
	Cfg := Config{}
	var cfgFile = "config.yml"
	cfgFileBytes, _ := ioutil.ReadFile(cfgFile)
	err := yaml.Unmarshal(cfgFileBytes, &Cfg)
	if err != nil {
		log.Print("unable to load config", err)
	}
	return Cfg
}
