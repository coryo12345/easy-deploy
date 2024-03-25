package docker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/coryo12345/easy-deploy/internal/config"
)

type DockerRepository interface {
	Health() (bool, error)
	GetStatuses(configEntries []config.ConfigEntry) ([]ConfigStatus, error)
	GetStatus(configEntry config.ConfigEntry) (DockerStatus, error)
	CloneRepo(config config.ConfigEntry) error
	BuildImage(config config.ConfigEntry) error
	StopContainer(config config.ConfigEntry) error
	DeleteContainer(config config.ConfigEntry) error
	StartContainer(config config.ConfigEntry) error
	CleanWorkDir(config config.ConfigEntry) error
}

type dockerRepo struct {
	workDir string
}

func New(workDir string) DockerRepository {
	return &dockerRepo{
		workDir: workDir,
	}
}

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

func (d dockerRepo) Health() (bool, error) {
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

func (d dockerRepo) GetStatuses(configEntries []config.ConfigEntry) ([]ConfigStatus, error) {
	statuses := make([]ConfigStatus, len(configEntries))
	for i, entry := range configEntries {
		status, err := d.GetStatus(entry)
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

func (d dockerRepo) GetStatus(configEntry config.ConfigEntry) (DockerStatus, error) {
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

func (d dockerRepo) CloneRepo(config config.ConfigEntry) error {
	return nil
}

func (d dockerRepo) BuildImage(config config.ConfigEntry) error {
	// docker build ...
	return nil
}

func (d dockerRepo) StopContainer(config config.ConfigEntry) error {
	// docker stop ...
	return nil
}

func (d dockerRepo) DeleteContainer(config config.ConfigEntry) error {
	// docker rm ...
	return nil
}

func (d dockerRepo) StartContainer(config config.ConfigEntry) error {
	// docker run ...
	return nil
}

func (d dockerRepo) CleanWorkDir(config config.ConfigEntry) error {
	return nil
}
