package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"main/internal/generic"
)

func main() {
	currentShell := console.GetCurrentShell()
	if !currentShell.IsValid() {
		fmt.Println(common.Unknown)
	} else {
		shells := generic.GetAllShellProfilesFiles()
		shell := shells[currentShell]
		if console.IsShell([]console.ShellType{console.PowerShell, console.Cmd}) {
			shell = shells[console.PowerShell]
		}
		if shell != "" {
			if !file.IsFile(shell) {
				file.WriteFile(shell, "", false, false)
			}
		} else {
			shell = common.Unknown
		}
		fmt.Println(shell)
	}
}
