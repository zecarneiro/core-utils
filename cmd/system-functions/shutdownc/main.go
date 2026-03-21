package main

import (
	"golangutils/pkg/logic"
	"golangutils/pkg/system"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "shutdownc",
		Short: "Shutdown the PC",
	}
	cobralib.WithRun(process)
}

func process() {
	logic.ProcessError(system.Shutdown())
}

func main() {
	cobralib.Run()
}
