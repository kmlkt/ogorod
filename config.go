package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sites []Site
}

var configPath = "./settings.yaml"

func ReadConfig() Config {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	StupidHandle(err)
	return config
}

func SaveConfig(config Config) {
	file, err := yaml.Marshal(config)
	StupidHandle(err)
	err = os.WriteFile(configPath, file, 0666)
	StupidHandle(err)
}

func (c Config) Apply() {
	c.Download()
	c.ConfigureNginx()
	RestartNginx()
}
