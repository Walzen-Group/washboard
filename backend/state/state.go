package state

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kpango/glg"
	"github.com/patrickmn/go-cache"
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
		instance.Config = Config{}
		instance.StackUpdateQueue = cache.New(5*time.Minute, 10*time.Minute)
		instance.StateQueue = cache.New(1*time.Minute, 1*time.Minute)
		reflectionPath = filepath.Dir(ex)
		fmt.Println(ex)
		// check if secrets file exists
		if _, err := os.Stat(filepath.Join(reflectionPath, "secrets.yaml")); os.IsNotExist(err) {
			encoded, err := yaml.Marshal(instance.Config)
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

		err = yaml.Unmarshal(yamlFile, &instance.Config)
		if err != nil {
			panic(err)
		}
	})
	return instance
}

func ReflectionPath() string {
	return reflectionPath
}

type Config struct {
	// secrets
	PortainerSecret string `yaml:"portainer_secret"`
	PortainerUrl    string `yaml:"portainer_url"`
	DbUrl           string `yaml:"db_url"`
}

type Data struct {
	Config           Config
	StackUpdateQueue *cache.Cache
	StateQueue *cache.Cache
}
