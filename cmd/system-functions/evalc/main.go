package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	command string
	args    []string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "evalc",
		Short: "Run given commands",
	}
	cobralib.CobraCmd.Flags().StringVarP(&command, "command", "c", "", "Command to run (required)")
	cobralib.CobraCmd.Flags().StringArrayVarP(&args, "arguments", "a", []string{}, "Arguments for command")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("command"))
	cobralib.WithRun(process)
}

func process() {
	if str.IsEmpty(command) {
		logic.ProcessError(fmt.Errorf("Invalid given command!"))
	}
	exe.ExecRealTime(models.Command{Cmd: command, Args: args, Verbose: true, UseShell: true})
}

func main() {
	cobralib.Run()
}
