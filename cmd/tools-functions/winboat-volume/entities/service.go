package entities

type Service struct {
	Image         string            `yaml:"image,omitempty"`
	ContainerName string            `yaml:"container_name,omitempty"`
	Environment   map[string]string `yaml:"environment,omitempty"`
	CapAdd        []string          `yaml:"cap_add,omitempty"`
	Privileged    bool              `yaml:"privileged,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	StopGrace     string            `yaml:"stop_grace_period,omitempty"`
	Restart       string            `yaml:"restart,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	Devices       []string          `yaml:"devices,omitempty"`
}
