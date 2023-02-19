package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/gommon/log"
	"sync"
)

type Config struct {
	ConfigFile string
	Listen     struct {
		Type   string `yaml:"type" default:"port"`
		BindIP string `yaml:"bind_ip" default:"127.0.0.1"`
		Port   string `yaml:"port" default:"8081"`
	} `yaml:"listen"`
	Store StoreConfig `yaml:"storage"`
}

type StoreConfig struct {
	Host     string `required:"true"`
	Port     int    `required:"true"`
	Database string `required:"true"`
	Username string `required:"true"`
	Password string `required:"true"`
}

var instance *Config
var once sync.Once

func NewConfig() *Config {

	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Info(help)
			log.Fatal(err)
		}
	})
	return instance
}
