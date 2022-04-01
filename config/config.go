package config

import (
	"github.com/go-yaml/yaml"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func GetConfig() (*Config, error) {

	newConfig := &Config{}

	file, err := os.Open("config/config.yml")
	if err != nil {
		log.Panicln("Could not find config file", err)
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&newConfig); err != nil {
		log.Panicln("Could not parse config.yml", err)
		return nil, err
	}

	return newConfig, nil
}
