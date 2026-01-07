package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/shell"

	"main/internal/libs"
)

func main() {
	currentShell := shell.GetCurrentShell()
	if !currentShell.IsValid() {
		fmt.Println(common.Unknown)
	} else {
		shells := libs.GetAllShellProfilesFiles()
		shellType := shells[currentShell]
		if shell.IsShell([]shell.ShellType{shell.PowerShell, shell.Cmd}) {
			shellType = shells[shell.PowerShell]
		}
		if shellType != "" {
			if !file.IsFile(shellType) {
				logic.ProcessError(file.WriteFile(shellType, "", false, false))
			}
		} else {
			shellType = common.Unknown
		}
		fmt.Println(shellType)
	}
}
