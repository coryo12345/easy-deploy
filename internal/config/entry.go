package config

type ConfigEntry struct {
	Id               string `json:"id"`
	GitRepository    string `json:"repo"`
	ContainerName    string `json:"container_name"`
	ContainerOptions string `json:"container_options"`
}
