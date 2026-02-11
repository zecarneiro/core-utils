package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/str"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "cmd-exists <command>",
		Short: "Check if command exists",
		Args:  cobra.ExactArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(command string) {
	if str.IsEmpty(command) {
		logic.ProcessError(fmt.Errorf("Invalid given command!"))
	}
	if platform.IsWindows() || shell.IsShell([]enums.ShellType{enums.PowerShell, enums.Cmd}) {
		cmd := fmt.Sprintf("if (Get-Command %s -ErrorAction SilentlyContinue) { Write-Output true } else { Write-Output false }", command)
		logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmd, UseShell: true}))
	} else if platform.IsLinux() {
		cmd := fmt.Sprintf("if command -v %s >/dev/null 2>&1; then echo true; else echo false; fi", command)
		logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmd, UseShell: true}))
	} else {
		fmt.Println(common.Unknown)
	}
}

func main() {
	cobralib.Run()
}
