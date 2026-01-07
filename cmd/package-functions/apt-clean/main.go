package main

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "apt-clean",
		Short: "Cleanup APT",
	}
	cobralib.WithRun(process)
}

func process() {
	aptCmdList := []string{
		"sudo apt clean -y",
		"sudo apt autoremove -y",
		"sudo apt autopurge -y",
		"sudo apt autoclean -y",
	}
	logger.Title("Cleanup APT")
	for _, aptCmd := range aptCmdList {
		logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: aptCmd, Verbose: true, UseShell: true}))
	}
}

func main() {
	cobralib.Run()
}
