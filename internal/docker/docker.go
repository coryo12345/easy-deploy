package docker

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/coryo12345/easy-deploy/internal/config"
)

type DockerStatus struct {
	ContainerName string `json:"container_name"`
	ImageName     string `json:"image_name"`
	Status        string `json:"status"`
}

func Health() (bool, error) {
	cmd := exec.Command("docker", "ps")
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()
	if out.String() != "" && err.String() == "" {
		return true, nil
	} else {
		return false, errors.New(err.String())
	}
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
	// docker run ...
	return nil
}

func BuildImage() error {
	// docker build ...
	return nil
}

func DeleteContainer() error {
	// docker rm ...
	return nil
}

func StopContainer() error {
	// docker stop ...
	return nil
}
