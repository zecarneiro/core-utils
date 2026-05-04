package models

type MJetBrains struct {
	Settings map[string]any `json:"settings,omitempty"`
	Plugins  []string       `json:"plugins,omitempty"`
}
