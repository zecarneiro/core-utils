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
		Use:   "scoop-clean",
		Short: "Cleanup SCOOP",
	}
	cobralib.WithRun(process)
}

func process() {
	aptCmdList := []string{
		"scoop cleanup --all",
		"scoop cache rm *",
	}
	logger.Title("Cleanup SCOOP")
	for _, aptCmd := range aptCmdList {
		logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: aptCmd, Verbose: true, UseShell: true}))
	}
}

func main() {
	cobralib.Run()
}
