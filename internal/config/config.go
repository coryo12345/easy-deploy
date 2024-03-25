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
}

type ConfigData struct {
	Init     string        `json:"init"`
	Services []ConfigEntry `json:"services"`
}

func New(configFile string) (ConfigRepository, error) {
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
	var configData ConfigData
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse file %s, it may not be valid JSON", configFile)
	}

	return configData, nil
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
