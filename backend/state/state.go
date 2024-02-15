package state

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kpango/glg"
	"gopkg.in/yaml.v2"
)

var instance *Data
var once sync.Once
var reflectionPath string

func Instance() *Data {
	once.Do(func() {
		ex, err := os.Executable()
		if err != nil {
			glg.Fatal(err)
		}
		instance = new(Data)
		reflectionPath = filepath.Dir(ex)
		fmt.Println(ex)
		// check if secrets file exists
		if _, err := os.Stat(filepath.Join(reflectionPath, "secrets.yaml")); os.IsNotExist(err) {
			encoded, err := yaml.Marshal(instance)
			if err != nil {
				glg.Fatal("could not marshal base config file")
			}
			if err := os.WriteFile(filepath.Join(reflectionPath, "secrets.yaml"), encoded, 0644); err != nil {
				glg.Fatal("could not save base config file")
			}
		}
		// read secrets from yaml file
		yamlFile, err := os.ReadFile(filepath.Join(reflectionPath, "secrets.yaml"))
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			panic(err)
		}
	})
	return instance
}

func ReflectionPath() string {
	return reflectionPath
}

type Data struct {
	// secrets
	PortainerSecret string `yaml:"portainer_secret"`
	PortainerUrl string `yaml:"portainer_url"`
}
