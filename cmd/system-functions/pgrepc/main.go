package main

import (
	"errors"
	"golangutils/pkg/console"
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
		Use:   "pgrepc <pattern>",
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
	if platform.IsLinux() {
		pgrepCmd, err := console.Which("pgrep")
		logic.ProcessError(err)
		cmd.Cmd = pgrepCmd
		cmd.Args = []string{"-l", patternArg}
		cmd.UseShell = false
	} else if platform.IsWindows() {
		pgrepcScript := libs.GetScriptCmdPathByName("pgrepc.ps1", "system")
		cmd.Cmd = pgrepcScript
		cmd.Args = []string{"-Pattern", patternArg}
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
