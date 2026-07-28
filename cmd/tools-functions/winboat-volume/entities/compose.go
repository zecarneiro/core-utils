package entities

type Compose struct {
	Name     string             `yaml:"name,omitempty"`
	Volumes  map[string]any     `yaml:"volumes,omitempty"`
	Services map[string]Service `yaml:"services,omitempty"`
}
