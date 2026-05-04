package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

type VSCodeUserDataProfiles struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type VSCodeStorageInfo struct {
	UserDataProfiles []VSCodeUserDataProfiles `json:"userDataProfiles"`
}

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "vscode-extension-upgrade",
		Short: "Update VSCode Extensions",
	}
	cobralib.WithRun(process)
}

func process() {
	var settingsDir string
	if platform.IsWindows() {
		settingsDir = file.JoinPath(system.HomeDir(), "AppData\\Roaming\\Code\\User")
	} else if platform.IsLinux() {
		settingsDir = file.JoinPath(system.HomeDir(), ".config/Code/User")
	}
	vscodeStorageInfo, err := file.ReadJsonFile[VSCodeStorageInfo](file.JoinPath(settingsDir, "globalStorage", "storage.json"))
	if err != nil {
		logger.Error(err)
	} else {
		cmd := `code --profile "%s" --update-extensions`
		logger.Title("Update All Extensions from all VSCode Profiles")
		for _, userDataProfile := range vscodeStorageInfo.UserDataProfiles {
			exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf(cmd, userDataProfile.Name), Verbose: true, UseShell: true})
		}
	}
}

func main() {
	cobralib.Run()
}
