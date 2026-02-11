package main

import (
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/obj"
	"golangutils/pkg/system"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "is-admin",
		Short: "Check if the app is running as admin or not",
	}
	cobralib.WithRun(process)
}

func process() {
	confirm, err := obj.ObjectToString(system.IsAdmin())
	logic.ProcessError(err)
	logger.Log(confirm)
}

func main() {
	cobralib.Run()
}
