package cfg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	t.Run("given no config.yml when called then it will still return valid default configuration settings", func(t *testing.T) {
		cwd, _ := os.Getwd()
		tmpDir, _ := ioutil.TempDir("", "dirstat_config_test")
		_ = os.Chdir(tmpDir)

		cwd2, _ := os.Getwd()
		fmt.Println(cwd2)

		config := GetConfig("")

		if config.ServicePort != "9999" {
			t.Fail()
			t.Errorf("the configuration returned must define the default service port 9999")
		}

		_ = os.Chdir(cwd)
	})

	t.Run("given config.yml is in executable directory when called then it will return config from yml file", func(t *testing.T) {
		cwd, _ := os.Getwd()
		tmpDir, _ := ioutil.TempDir("", "dirstat_config_test")
		_ = os.Chdir(tmpDir)

		cwd2, _ := os.Getwd()
		fmt.Println(cwd2)

		configSrc := Config{
			ServicePort: "9997",
		}

		out, _ := yaml.Marshal(configSrc)

		_ = ioutil.WriteFile("config.yml", []byte(out), 0644)

		config := GetConfig("")

		if config.ServicePort != "9997" {
			t.Fail()
			t.Errorf("the configuration returned must return the service port 9997, defined in config.yml")
			t.Errorf(" it returned %v instead.", config.ServicePort)
		}

		_ = os.Chdir(cwd)
	})
	t.Run("given a custom config.yml not in the current directory when called with its path as argument then it will return config from custom yml file", func(t *testing.T) {
		cwd, _ := os.Getwd()
		tmpDir, _ := ioutil.TempDir("", "dirstat_config_test")
		_ = os.Chdir(tmpDir)

		cwd2, _ := os.Getwd()
		fmt.Println(cwd2)

		configSrc := Config{
			ServicePort: "9994",
		}

		out, _ := yaml.Marshal(configSrc)

		customConfigFile := "myconfig.yml"
		_ = ioutil.WriteFile(customConfigFile, []byte(out), 0644)

		config := GetConfig(customConfigFile)

		if config.ServicePort != "9994" {
			t.Fail()
			t.Errorf("the configuration returned must return the service port 9997, defined in config.yml")
			t.Errorf(" it returned %v instead.", config.ServicePort)
		}

		_ = os.Chdir(cwd)
	})
}
