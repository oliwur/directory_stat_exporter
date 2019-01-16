package cfg

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Dir struct {
	Name      string
	Path      string
	Recursive bool
}
type Config struct {
	ServicePort string
	CacheTime   int
	Directories []Dir
}

func GetConfig(fileName string) Config {
	Cfg := Config{}
	// set default values
	Cfg.ServicePort = "9999"
	Cfg.CacheTime = 5

	var cfgFile = "config.yml"
	if fileName != "" {
		cfgFile = fileName
	}
	cfgFileBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Print("unable to load config", err)
	} else {
		err = yaml.Unmarshal(cfgFileBytes, &Cfg)
		if err != nil {
			log.Print("unable to load config", err)
		}
	}

	return Cfg
}
