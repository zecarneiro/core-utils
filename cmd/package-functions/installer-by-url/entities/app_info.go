package entities

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/system"

	"main/cmd/package-functions/installer-by-url/args"
)

type AppInfo struct {
	Name       string         `json:"name,omitempty"`
	Url        string         `json:"url,omitempty"`
	Extension  string         `json:"extension,omitempty"`
	Source     string         `json:"source,omitempty"`
	Github     GithubInfo     `json:"github"`
	FileRunner FileRunnerData `json:"fileRunner"`
}

func NewAppInfo(baseArgs *args.BaseArgs, destScriptDir string, source string) *AppInfo {
	logger.Header(fmt.Sprintf("Init install %s", baseArgs.AppName))
	savedScriptFile := ""
	commandStr := ""
	if !baseArgs.IsCommand {
		fileBasename := file.Basename(baseArgs.ScriptOrCommand)
		savedScriptFile = file.JoinPath(destScriptDir, fileBasename)
	} else {
		commandStr = baseArgs.ScriptOrCommand
	}
	return &AppInfo{
		Name:      baseArgs.AppName,
		Extension: baseArgs.FileExtension,
		Url:       baseArgs.Url,
		Source:    source,
		FileRunner: FileRunnerData{
			IsCommand:   baseArgs.IsCommand,
			ShellScript: savedScriptFile,
			Command:     commandStr,
			Shell:       baseArgs.Shell,
			OutputFile:  file.JoinPath(system.TempDir(), fmt.Sprintf(`%s.%s`, baseArgs.AppName, baseArgs.FileExtension)),
		},
	}
}
