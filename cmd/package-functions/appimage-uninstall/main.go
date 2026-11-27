package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/system"
	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	commandOnlyArg bool
	forceArg       bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "appimage-uninstall <app>",
		Short: "Uninstall appimage app",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().BoolVarP(&commandOnlyArg, "command-only", "c", false, "Uninstall command only.")
	cobralib.CobraCmd.Flags().BoolVarP(&forceArg, "force", "f", false, "Force uninstall.")
	cobralib.WithRunArgsStr(process)
}

func getInstalledApp(app string) *libs.GearlevelInstalledAppInfo {
	var findedApp *libs.GearlevelInstalledAppInfo
	for _, appInstalled := range libs.GetGearlevelListOfInstalledApps().Installed {
		if appInstalled.GetRegisterName() == app {
			findedApp = &appInstalled
			break
		}
	}
	if findedApp == nil && forceArg {
		appImageDir := fmt.Sprintf(`%s/AppImages`, system.HomeDir())
		dirInfo, err := file.ReadDir(appImageDir)
		findedApp = &libs.GearlevelInstalledAppInfo{
			Id:           fmt.Sprintf("%s.desktop", app),
			AppimageFile: "",
		}
		logic.ProcessError(err)
		for _, fileDir := range dirInfo.Files {
			if file.FileName(fileDir) == app {
				findedApp.AppimageFile = fmt.Sprintf(`%s/%s`, appImageDir, fileDir)
				break
			}
		}
	}
	return findedApp
}

func process(app string) {
	installedApp := getInstalledApp(app)
	if installedApp != nil {
		registerName := installedApp.GetRegisterName()
		if !commandOnlyArg {
			cmdStr := fmt.Sprintf("flatpak run it.mijorus.gearlever --remove '%s' -y", installedApp.AppimageFile)
			desktopFile := installedApp.GetDesktopFile()
			err := exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: false})
			if !forceArg {
				logic.ProcessError(err)
			} else {
				logger.Error(err)
			}
			if file.IsFile(desktopFile) {
				logger.Info(fmt.Sprintf("Deleting .desktop file: %s", desktopFile))
				logger.Error(file.DeleteFile(desktopFile))
			}
			if file.IsFile(installedApp.AppimageFile) {
				logger.Info(fmt.Sprintf("Deleting .appimage file: %s", installedApp.AppimageFile))
				logger.Error(file.DeleteFile(installedApp.AppimageFile))
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
