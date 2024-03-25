package config

type ConfigEntry struct {
	Id               string            `json:"id"`
	GitRepository    string            `json:"repo"`
	ContainerName    string            `json:"container_name"`
	ContainerOptions string            `json:"container_options"`
	ImageName        string            `json:"image_name"`
	Dockerfile       string            `json:"dockerfile_path"`
	Env              map[string]string `json:"env"`
}
