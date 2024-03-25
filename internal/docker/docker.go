package docker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/coryo12345/easy-deploy/internal/config"
)

type ConfigStatus struct {
	Config config.ConfigEntry
	Status DockerStatus
	Error  error
}

type DockerStatus struct {
	Command      string
	CreatedAt    string
	ID           string
	Image        string
	Labels       string
	LocalVolumes string
	Mounts       string
	Names        string
	Networks     string
	Ports        string
	RunningFor   string
	Size         string
	State        string
	Status       string
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

func GetStatuses(configEntries []config.ConfigEntry) ([]ConfigStatus, error) {
	statuses := make([]ConfigStatus, len(configEntries))
	for i, entry := range configEntries {
		status, err := GetStatus(entry)
		if err != nil {
			statuses[i] = ConfigStatus{
				Config: entry,
				Status: DockerStatus{},
				Error:  err,
			}
		} else {
			statuses[i] = ConfigStatus{
				Config: entry,
				Status: status,
				Error:  nil,
			}
		}
	}
	return statuses, nil
}

func GetStatus(configEntry config.ConfigEntry) (DockerStatus, error) {
	cmd := exec.Command("docker", "ps", "--format={{json .}}", "-f", fmt.Sprintf("name=%s", configEntry.ContainerName))
	out, err := cmd.Output()
	if err != nil {
		return DockerStatus{}, err
	}

	status := DockerStatus{}
	err = json.Unmarshal(out, &status)
	if err != nil {
		return DockerStatus{}, err
	}
	return status, nil

}

func CloneRepo(workdir string, repo string) error {
	return nil
}

func BuildImage(workdir string, config config.ConfigEntry) error {
	// docker build ...
	return nil
}

func StopContainer(config config.ConfigEntry) error {
	// docker stop ...
	return nil
}

func DeleteContainer(config config.ConfigEntry) error {
	// docker rm ...
	return nil
}

func StartContainer(config config.ConfigEntry) error {
	// docker run ...
	return nil
}

func CleanWorkDir(workdir string, config config.ConfigEntry) error {
	return nil
}
