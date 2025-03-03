package sys

import (
	"encoding/json"
	"log/slog"
	"os"
	"slices"
	"sync"
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

func (c Config) DoStuff(pm *ProcessMap, changed []ShortID, register func(ProcessMap)) error {
	firstTime := false
	if *pm == nil {
		*pm = make(ProcessMap)
		firstTime = true
	}
	newPm, err := c.apply(firstTime, changed)
	if err != nil {
		return err
	}
	newPm = Fill(*pm, newPm)
	register(newPm)
	CollectGarbage(*pm, newPm, false)
	*pm = newPm
	return nil
}

func (c Config) apply(firstTime bool, changed []ShortID) (ProcessMap, error) {
	pm := make(ProcessMap)
	wg := sync.WaitGroup{}
	wg.Add(len(c.Services))
	for _, service := range c.Services {
		go func() {
			process, err := service.DoStuff(firstTime || slices.Contains(changed, service.ID),
				slices.Contains(changed, service.ID))
			if err != nil {
				slog.Error("Service failed", "service", service.ID, "error", err)
				wg.Done()
				return
			}
			slog.Debug("Service done", "service", service.ID, "process", process)

			pm[service.ID] = process
			wg.Done()
		}()
	}
	wg.Wait()
	return pm, nil
}
