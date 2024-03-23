package docker

import (
	"errors"

	"github.com/coryo12345/easy-deploy/internal/config"
)

type DockerStatus struct {
	ContainerName string `json:"container_name"`
	ImageName     string `json:"image_name"`
	Status        string `json:"status"`
}

func Health() (bool, error) {
	// make sure docker can be run
	return false, errors.New("TODO")
}

func GetStatuses(configEntries []config.ConfigEntry) ([]DockerStatus, error) {
	statuses := make([]DockerStatus, len(configEntries))
	for i, entry := range configEntries {
		status, err := GetStatus(entry)
		if err != nil {
			return nil, err
		}
		statuses[i] = status
	}
	return statuses, nil
}

func GetStatus(configEntry config.ConfigEntry) (DockerStatus, error) {
	return DockerStatus{
		ContainerName: configEntry.ContainerName,
		ImageName:     configEntry.ImageName,
		Status:        "TODO",
	}, nil
}

func StartContainer() error {
	return nil
}

func BuildImage() error {
	return nil
}

func DeleteContainer() error {
	return nil
}

func StopContainer() error {
	return nil
}
