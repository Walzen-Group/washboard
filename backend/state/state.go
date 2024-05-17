package state

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
		instance.Config = Config{CacheDurationMinutes: 1, StartStacksOnLaunch: false, StartEndpointId: 1}
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

		// overwrite config with environment variables if they exist
		overwriteConfigWithEnv(&instance.Config)

		// save config file with all fields
		encoded, err := yaml.Marshal(instance.Config)
		if err != nil {
			glg.Fatal("could not marshal base config file")
		}
		if err := os.WriteFile(filepath.Join(reflectionPath, "secrets.yaml"), encoded, 0644); err != nil {
			glg.Fatal("could not save base config file")
		}

		if instance.Config.CacheDurationMinutes == 1 {
			glg.Warn("image cache duration is set to 1 minute")
		} else {
			glg.Infof("image cache duration is set to %d minutes", instance.Config.CacheDurationMinutes)
		}
	})
	return instance
}

func ReflectionPath() string {
	return reflectionPath
}

type Config struct {
	// secrets
	PortainerSecret      string   `yaml:"portainer_secret"`
	PortainerUrl         string   `yaml:"portainer_url"`
	DbUrl                string   `yaml:"db_url"`
	User                 string   `yaml:"user"`
	Password             string   `yaml:"password"`
	JwtSecret            string   `yaml:"jwt_secret"`
	CacheDurationMinutes int      `yaml:"cache_duration_minutes"`
	Cors                 []string `yaml:"cors"`
	StartStacksOnLaunch  bool     `yaml:"start_stacks_on_launch"`
	StartEndpointId      int      `yaml:"start_endpoint_id"`
}

type Data struct {
	Config           Config
	StackUpdateQueue *cache.Cache
	StateQueue       *cache.Cache
}

func overwriteConfigWithEnv(config *Config) {
	if value, exists := os.LookupEnv("PORTAINER_SECRET"); exists {
		config.PortainerSecret = value
	}
	if value, exists := os.LookupEnv("PORTAINER_URL"); exists {
		config.PortainerUrl = value
	}
	if value, exists := os.LookupEnv("DB_URL"); exists {
		config.DbUrl = value
	}
	if value, exists := os.LookupEnv("USER"); exists {
		config.User = value
	}
	if value, exists := os.LookupEnv("PASSWORD"); exists {
		config.Password = value
	}
	if value, exists := os.LookupEnv("JWT_SECRET"); exists {
		config.JwtSecret = value
	}
	if value, exists := os.LookupEnv("CACHE_DURATION_MINUTES"); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			config.CacheDurationMinutes = intValue
		} else {
			glg.Warn("invalid CACHE_DURATION_MINUTES value, using default")
		}
	}
	if value, exists := os.LookupEnv("CORS"); exists {
		cors := strings.Split(value, ",")
		for i := range cors {
			cors[i] = strings.TrimSpace(cors[i])
		}
		config.Cors = cors
	}
	if value, exists := os.LookupEnv("START_STACKS_ON_LAUNCH"); exists {
		if value == "true" {
			config.StartStacksOnLaunch = true
		} else {
			config.StartStacksOnLaunch = false
		}
	}
}
