package main

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() {
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "deb-get-clean",
		Short: "Cleanup DEB-GET",
	}
	cobralib.WithRun(process)
}

func process() {
	logger.Title("Cleanup DEB-GET")
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: "sudo deb-get clean", Verbose: true, UseShell: true}))
}

func main() {
	cobralib.Run()
}
