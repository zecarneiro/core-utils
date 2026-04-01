package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
)

var (
	command string
	args    []string
)

func init() { setupCommand() }

func setupCommand() {
	args = console.GetArgsList()
}

func process() {
	cmdArgs := slice.ArrayToString(args)
	if str.IsEmpty(cmdArgs) {
		logic.ProcessError(fmt.Errorf("Invalid given command!"))
	}
	exe.ExecRealTime(models.Command{Cmd: cmdArgs, Verbose: true, UseShell: true})
}

func main() {
	process()
}
