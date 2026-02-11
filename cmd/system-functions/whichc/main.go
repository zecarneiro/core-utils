package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/shell"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var result string

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "whichc <command>",
		Short: "Find given command",
		Args:  cobra.ExactArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func runCmd(cmd models.Command) string {
	res, _ := exe.Exec(cmd)
	return res
}

func findCmdPath(command string) bool {
	path, _ := console.Which(command)
	result = path
	exist := logic.Ternary(len(path) > 0, true, false)
	if exist {
		result = fmt.Sprintf("APP/CMD: %s", path)
	}
	return exist
}

func findAlias(command string) bool {
	exist := false
	canRun := true
	cmd := models.Command{UseShell: true}
	if shell.IsPowerShellSimple() {
		cmd.Cmd = fmt.Sprintf("Get-Alias %s -ErrorAction SilentlyContinue", command)
	} else if shell.IsBashSimple() {
		cmd.Cmd = fmt.Sprintf("alias %s >/dev/null 2>&1 && echo existe", command)
		cmd.IsInteractiveShell = true
	} else if shell.IsFishSimple() {
		cmd.Cmd = fmt.Sprintf("alias | grep -q '^alias %s '; and echo exists", command)
	} else if shell.IsKshSimple() || shell.IsZshSimple() {
		cmd.Cmd = fmt.Sprintf("alias | grep -q \"^%s=\" && echo exists", command)
		cmd.IsInteractiveShell = true
	} else {
		canRun = false
	}
	if canRun {
		path := runCmd(cmd)
		if len(path) > 0 {
			exist = true
		}
		if exist {
			result = fmt.Sprintf("ALIAS: %s", command)
		}
	}
	return exist
}

func findFuncion(command string) bool {
	exist := false
	canRun := true
	cmd := models.Command{UseShell: true}
	if shell.IsPowerShellSimple() {
		cmd.Cmd = fmt.Sprintf("Get-ChildItem Function:%s -ErrorAction SilentlyContinue", command)
	} else if shell.IsBashSimple() {
		cmd.Cmd = fmt.Sprintf("declare -F %s >/dev/null 2>&1 && echo existe", command)
		cmd.IsInteractiveShell = true
	} else if shell.IsFishSimple() {
		cmd.Cmd = fmt.Sprintf("functions -q %s; and echo \"exists\"", command)
	} else if shell.IsKshSimple() || shell.IsZshSimple() {
		cmd.Cmd = fmt.Sprintf("typeset -f %s >/dev/null 2>&1 && echo exists", command)
		cmd.IsInteractiveShell = true
	} else {
		canRun = false
	}
	if canRun {
		path := runCmd(cmd)
		if len(path) > 0 {
			exist = true
		}
		if exist {
			result = fmt.Sprintf("FUNCTION: %s", command)
		}
	}
	return exist
}

func process(command string) {
	if str.IsEmpty(command) {
		logic.ProcessError(fmt.Errorf("Invalid given command!"))
	}
	if findCmdPath(command) || findAlias(command) || findFuncion(command) {
		logger.Log(result)
	} else {
		logger.Log(common.Unknown)
	}
}

func main() {
	cobralib.Run()
}
