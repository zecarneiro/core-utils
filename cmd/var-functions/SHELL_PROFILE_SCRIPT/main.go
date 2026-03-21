package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/shell"
)

func main() {
	currentShell := shell.GetCurrentShellSimple()
	if !currentShell.IsValid() {
		fmt.Println(common.Unknown)
	} else {
		var shellFile string
		if shell.IsShell([]enums.ShellType{enums.PowerShell, enums.Cmd}) {
			shellFile = shell.GetShellProfileFile(enums.PowerShell)
		} else {
			shellFile = shell.GetShellProfileFile(currentShell)
		}
		if shellFile != common.Unknown && !file.IsFile(shellFile) {
			fileConfig := models.FileWriterConfig{
				File:        shellFile,
				Data:        "",
				IsAppend:    false,
				IsCreateDir: false,
			}
			logic.ProcessError(file.WriteFile(fileConfig))
		}
		fmt.Println(shellFile)
	}
}
