package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		DB    Postgres `yaml:"db"`
		Redis Redis    `yaml:"redis"`
	}
	Postgres struct {
		URL        string `yaml:"url"`
		DriverName string `yaml:"driver_name"`
	}
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
)

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	file, err := os.ReadFile("./config/config.yml")
	if err != nil {
		log.Errorf("Failed to open config file: %v", err)
		return cfg, err
	}

	// Parse the YAML into the Config struct
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Errorf("Failed to parse config file: %v", err)
		return cfg, err
	}

	return cfg, nil

}

func LoadTestConfig() (*Config, error) {
	cfg := &Config{}

	file, err := os.ReadFile("../config/config_test.yml")
	if err != nil {
		log.Errorf("Failed to open config file: %v", err)
		return cfg, err
	}

	// Parse the YAML into the Config struct
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Errorf("Failed to parse config file: %v", err)
		return cfg, err
	}

	return cfg, nil

}
