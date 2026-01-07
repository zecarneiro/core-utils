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
		Use:   "flatpak-clean",
		Short: "Cleanup Flatpak",
	}
	cobralib.WithRun(process)
}

func process() {
	cmdList := []string{
		"flatpak uninstall --unused -y",
		"sudo rm -rfv /var/tmp/flatpak-cache*",
	}
	logger.Title("Cleanup Flatpak")
	for _, aptCmd := range cmdList {
		logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: aptCmd, Verbose: true, UseShell: true}))
	}
}

func main() {
	cobralib.Run()
}
