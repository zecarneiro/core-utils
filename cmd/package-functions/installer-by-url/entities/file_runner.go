package entities

type FileRunnerData struct {
	IsCommand   bool   `json:"isCommand,omitempty"`
	Command     string `json:"command,omitempty"`
	ShellScript string `json:"shellScript,omitempty"`
	Shell       string `json:"shell,omitempty"`
	OutputFile  string `json:"outputFile,omitempty"`
}
