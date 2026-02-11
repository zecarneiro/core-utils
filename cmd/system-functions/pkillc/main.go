package main

import (
	"errors"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "pkillc <pattern>",
		Short: "",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(patternArg string) {
	if str.IsEmpty(patternArg) {
		logic.ProcessError(errors.New("Invalid given Pattern"))
	}
	cmd := models.Command{
		UseShell: true,
	}
	if platform.IsWindows() {
		pgrepcScript := libs.GetScriptCmdPathByName("pkillc.ps1", "system")
		cmd.Cmd = pgrepcScript
		cmd.Args = []string{"-Pattern", patternArg}
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
