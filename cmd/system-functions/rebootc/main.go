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
		Use:   "rebootc",
		Short: "Reboot the PC",
	}
	cobralib.WithRun(process)
}

func process() {
	logic.ProcessError(system.Reboot())
}

func main() {
	cobralib.Run()
}
