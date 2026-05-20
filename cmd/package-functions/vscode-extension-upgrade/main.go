package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "vscode-extension-upgrade",
		Short: "Update VSCode Extensions",
	}
	cobralib.WithRun(process)
}

func process() {
	cmd := `code --profile "%s" --update-extensions`
	for _, userDataProfile := range libs.GetVscodeUserDataProfiles() {
		err := exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf(cmd, userDataProfile.Name), Verbose: true, UseShell: true})
		logger.Error(err)
	}
}

func main() {
	cobralib.Run()
}
