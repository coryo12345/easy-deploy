package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ConfigRepository interface {
	GetAllServices() []ConfigEntry
	FindEntryById(id string) (ConfigEntry, error)
	Refresh() error
}

type ConfigData struct {
	configFile string
	Init       string        `json:"init"`
	Services   []ConfigEntry `json:"services"`
}

func New(configFile string) (ConfigRepository, error) {
	data, err := readDataFromFile(configFile)
	if err != nil {
		return nil, err
	}

	var configData ConfigData
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse file %s, it may not be valid JSON", configFile)
	}

	configData.configFile = configFile

	return &configData, nil
}

func readDataFromFile(configFile string) ([]byte, error) {
	// open config file
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s", configFile)
	}
	defer file.Close()

	// parse to struct
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s", configFile)
	}
	return data, nil
}

func (c ConfigData) GetAllServices() []ConfigEntry {
	return c.Services
}

func (c ConfigData) FindEntryById(id string) (ConfigEntry, error) {
	for _, entry := range c.Services {
		if entry.Id == id {
			return entry, nil
		}
	}
	return ConfigEntry{}, fmt.Errorf("no entry with id %s found", id)
}

func (c *ConfigData) Refresh() error {
	data, err := readDataFromFile(c.configFile)
	if err != nil {
		return err
	}

	var newConfig ConfigData
	err = json.Unmarshal(data, &newConfig)
	if err != nil {
		return err
	}

	c.Init = newConfig.Init
	c.Services = newConfig.Services

	return nil
}
