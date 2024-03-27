package docker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/coryo12345/easy-deploy/internal/config"
)

const (
	WORKDIR_SUB_NAME = "easydeploy"
)

type DockerRepository interface {
	Health() (bool, error)
	GetStatuses(configEntries []config.ConfigEntry) ([]ConfigStatus, error)
	GetStatus(configEntry config.ConfigEntry) (DockerStatus, error)
	CloneRepo(config config.ConfigEntry, out *strings.Builder) error
	BuildImage(config config.ConfigEntry, out *strings.Builder) error
	StopContainer(config config.ConfigEntry, out *strings.Builder) error
	DeleteContainer(config config.ConfigEntry, out *strings.Builder) error
	StartContainer(config config.ConfigEntry, out *strings.Builder) error
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
	cmd := exec.Command("docker", "ps", "-a", "--format={{json .}}", "-f", fmt.Sprintf("name=%s", configEntry.ContainerName))
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

func (d dockerRepo) CloneRepo(config config.ConfigEntry, logs *strings.Builder) error {
	workDirPath := path.Join(d.workDir, WORKDIR_SUB_NAME, config.Id)
	err := os.MkdirAll(workDirPath, fs.ModePerm)
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "clone", config.GitRepository, workDirPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	logs.WriteString(out.String())
	if err != nil {
		return err
	}
	return nil
}

func (d dockerRepo) BuildImage(config config.ConfigEntry, logs *strings.Builder) error {
	workDirPath := path.Join(d.workDir, WORKDIR_SUB_NAME, config.Id)
	dockerfilePath := path.Join(workDirPath, config.Dockerfile)
	cmd := exec.Command("docker", "build", "-t", config.ImageName, "-f", dockerfilePath, workDirPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	logs.WriteString(out.String())
	if err != nil {
		return err
	}
	return nil
}

func (d dockerRepo) StopContainer(config config.ConfigEntry, logs *strings.Builder) error {
	cmd := exec.Command("docker", "stop", config.ContainerName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	logs.WriteString(out.String())
	if err != nil {
		// going to ignore this for now - if the container doesn't exist it will error
		log.Println(err.Error())
	}
	return nil
}

func (d dockerRepo) DeleteContainer(config config.ConfigEntry, logs *strings.Builder) error {
	cmd := exec.Command("docker", "rm", config.ContainerName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	logs.WriteString(out.String())
	if err != nil {
		// if the previous container didnt exist this would error
		log.Println(err.Error())
	}
	return nil
}

func (d dockerRepo) StartContainer(config config.ConfigEntry, logs *strings.Builder) error {
	envStr := strings.Builder{}
	for key, value := range config.Env {
		envStr.WriteString(fmt.Sprintf(" --env %s=%s", key, value))
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("docker run %s --name %s %s %s", envStr.String(), config.ContainerName, config.ContainerOptions, config.ImageName))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	logs.WriteString(out.String())
	if err != nil {
		return err
	}
	return nil
}

func (d dockerRepo) CleanWorkDir(config config.ConfigEntry) error {
	workDirPath := path.Join(d.workDir, WORKDIR_SUB_NAME, config.Id)
	err := os.RemoveAll(workDirPath)
	if err != nil {
		log.Printf("ERROR CLEANING WORKDIR FOR %s\n", config.Id)
		log.Println(err.Error())
	}
	return nil
}
