package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var commandOnlyArg bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "appimage-uninstall <app>",
		Short: "Uninstall appimage app",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().BoolVarP(&commandOnlyArg, "command-only", "c", false, "Uninstall command only.")
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	installedApps := libs.GetGearlevelListOfInstalledApps().Installed
	for _, appInstalled := range installedApps {
		registerName := appInstalled.GetRegisterName()
		if registerName == app {
			if !commandOnlyArg {
				cmdStr := fmt.Sprintf("flatpak run it.mijorus.gearlever --remove '%s' -y", appInstalled.AppimageFile)
				desktopFile := appInstalled.GetDesktopFile()
				logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: false}))
				if file.IsFile(desktopFile) {
					logger.Info(fmt.Sprintf("Deleting .desktop file: %s", desktopFile))
					logger.Error(file.DeleteFile(desktopFile))
				}
			}
			if console.CmdExists(registerName) {
				libs.RunCoreUtilsCmd("script-manager-cu", false, "uninstall", registerName)
			}
		}
	}
}

func main() {
	cobralib.Run()
}
