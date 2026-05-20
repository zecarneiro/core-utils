package entities

type GithubInfo struct {
	Owner            string `json:"owner,omitempty"`
	Repository       string `json:"repository,omitempty"`
	IsLatest         bool   `json:"latest,omitempty"`
	InstalledVersion string `json:"version,omitempty"`
}
