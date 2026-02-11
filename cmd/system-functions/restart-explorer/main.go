package main

import (
	"fmt"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"main/internal/libs/cobralib"
	"slices"

	"github.com/spf13/cobra"
)

var (
	typeFileExplorer = map[string]string{
		"nautilus":     "nautilus -q",
		"win-explorer": "Stop-Process -Name explorer -Force; Start-Process explorer.exe",
	}
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "restart-explorer [type_explorer]",
		Short: "Restart file explorer. If type_explorer not set, will be restart the system explorer(First founded by system)",
	}
	cobralib.WithRunArgsStr(process)
}

func process(typeExplorer string) {
	cmd := models.Command{UseShell: true}
	const errorMsg = "Can not restart the system explorer"
	canRun := true
	if !str.IsEmpty(typeExplorer) {
		typeExplorerArr := slice.MapToKeys(typeFileExplorer)
		if !slices.Contains(typeExplorerArr, typeExplorer) {
			logic.ProcessError(fmt.Errorf("type_explorer must be one of: %v", slice.ArrayToStringBySep(typeExplorerArr, ",")))
		}
	} else {
		osType := system.GetOsType()
		if !osType.IsValid() {
			canRun = false
		} else if osType.Equals(enums.UbuntuSO) {
			typeExplorer = "nautilus"
		} else if platform.IsWindows() {
			typeExplorer = "win-explorer"
		} else {
			canRun = false
		}
	}
	if canRun {
		cmd.Cmd = typeFileExplorer[typeExplorer]
		logic.ProcessError(exe.ExecRealTime(cmd))
	} else {
		logic.ProcessError(fmt.Errorf(errorMsg))
	}
}

func main() {
	cobralib.Run()
}
