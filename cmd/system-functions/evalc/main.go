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
		Use:   "evalc <command>",
		Short: "Run given commands",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(args string) {
	if str.IsEmpty(args) {
		logic.ProcessError(fmt.Errorf("Invalid given command!"))
	}
	exe.ExecRealTime(models.Command{Cmd: args, Verbose: true, UseShell: true})
}

func main() {
	cobralib.Run()
}
