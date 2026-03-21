package main

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var force bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "wsl-shutdown",
		Short: "Shutdown all WSL distro",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&force, "force", "f", false, "Shutdown force")
	cobralib.WithRun(process)
}

func process() {
	cmdInfo := models.Command{Cmd: "wsl --shutdown", UseShell: true}
	if force {
		cmdInfo.Cmd = "sudo taskkill /F /IM wslservice.exe"
	}
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
}

func main() {
	cobralib.Run()
}
