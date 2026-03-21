package main

import (
	"fmt"
	"golangutils/pkg/shell"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var deepVerify bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "SHELL_NAME",
		Short: "Show the current shell name",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&deepVerify, "deep-verify", "d", false, "Verify current shell deeper(More precise - Will take more time to check)")
	cobralib.WithRun(process)
}

func process() {
	if deepVerify {
		fmt.Println(shell.GetCurrentShell().String())
	} else {
		fmt.Println(shell.GetCurrentShellSimple().String())
	}
}

func main() {
	cobralib.Run()
}
