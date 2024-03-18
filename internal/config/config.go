package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ConfigRepository interface {
	GetAllEntries() []ConfigEntry
	FindEntryById(id string) (ConfigEntry, error)
}

type config struct {
	entries []ConfigEntry
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
	var entries []ConfigEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, fmt.Errorf("unable to parse file %s, it may not be valid JSON", configFile)
	}

	c := config{
		entries: entries,
	}

	return c, nil
}

func (c config) GetAllEntries() []ConfigEntry {
	return c.entries
}

func (c config) FindEntryById(id string) (ConfigEntry, error) {
	for _, entry := range c.entries {
		if entry.Id == id {
			return entry, nil
		}
	}
	return ConfigEntry{}, fmt.Errorf("no entry with id %s found", id)
}
