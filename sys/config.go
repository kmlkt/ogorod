package sys

import (
	"encoding/json"
	"os"
)

type Config struct {
	Services []Service
}

const ConfigPath = "./config.json"

func LoadConfig() (Config, error) {
	file, err := os.ReadFile(ConfigPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c Config) Save() error {
	file, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath, file, 0666)
}

func (c Config) DoStuff(pm *ProcessMap, register func(ProcessMap)) error {
	firstTime := false
	if *pm == nil {
		*pm = make(ProcessMap)
		firstTime = true
	}
	newPm, err := c.apply(firstTime)
	if err != nil {
		return err
	}
	newPm = Fill(*pm, newPm)
	register(newPm)
	CollectGarbage(*pm, newPm, false)
	*pm = newPm
	return nil
}

func (c Config) apply(firstTime bool) (ProcessMap, error) {
	pm := make(ProcessMap)
	for _, service := range c.Services {
		process, err := service.DoStuff(firstTime)
		if err != nil {
			return pm, err
		}

		pm[service.ID] = process
	}
	return pm, nil
}
