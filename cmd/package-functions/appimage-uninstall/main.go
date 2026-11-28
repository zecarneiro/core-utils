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

var (
	commandOnlyArg bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "appimage-uninstall <app>",
		Short: "Uninstall appimage app(By Gear Lever)",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().BoolVarP(&commandOnlyArg, "command-only", "c", false, "Uninstall command only.")
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	var foundedApp *libs.GearLeverInstalledAppInfo
	installedAppList := libs.GetGearLeverListOfInstalledApps().Installed
	for _, appInstalled := range installedAppList {
		registerName := appInstalled.GetRegisterName()
		if registerName == app {
			if !appInstalled.IsNotRegisted || console.Confirm("This app is not registed. Continue?", false) {
				foundedApp = &appInstalled
			}
			break
		}
	}
	if foundedApp != nil {
		registerName := foundedApp.GetRegisterName()
		if !commandOnlyArg {
			desktopFile := foundedApp.GetDesktopFile()
			if !foundedApp.IsNotRegisted {
				cmdStr := fmt.Sprintf("flatpak run it.mijorus.gearlever --remove '%s' -y", foundedApp.AppimageFile)
				err := exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: false})
				logic.ProcessError(err)
			}
			if file.IsFile(desktopFile) {
				logger.Info(fmt.Sprintf("Deleting .desktop file: %s", desktopFile))
				logger.Error(file.DeleteFile(desktopFile))
			}
			if file.IsFile(foundedApp.AppimageFile) {
				logger.Info(fmt.Sprintf("Deleting .appimage file: %s", foundedApp.AppimageFile))
				logger.Error(file.DeleteFile(foundedApp.AppimageFile))
			}
		}
		if console.CmdExists(registerName) {
			libs.RunCoreUtilsCmd("script-manager-cu", false, "uninstall", registerName)
		}
		logger.Ok("Uninstall done.")
	} else {
		logger.Error(fmt.Errorf(`Appimage app not found: %s`, app))
	}
}

func main() {
	cobralib.Run()
}
