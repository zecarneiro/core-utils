package main

import (
	"fmt"
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

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "appimage-install <package-file>",
		Short: "Install AppImage package.",
	}
	cobralib.WithRunArgsStr(process)
}

func process(packageFile string) {
	if !file.IsFile(packageFile) {
		logic.ProcessError(fmt.Errorf("Invalid given appimage file."))
	}
	packageFileTemp := fmt.Sprintf("%s/%s-copy.%s", system.TempDir(), file.FileName(packageFile), file.FileExtension(packageFile))
	logger.Info(fmt.Sprintf("Copying to '%s'", packageFileTemp))
	file.CopyFile(packageFile, packageFileTemp)
	logger.Info("Installing...")
	cmdStr := fmt.Sprintf("flatpak run it.mijorus.gearlever --integrate '%s' -y", packageFileTemp)
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: false}))
	installedApps := libs.GetGearlevelListOfInstalledApps().Installed
	for _, appInstalled := range installedApps {
		execCmdStr, err := system.GetExecDataFromDesktopFile(appInstalled.GetDesktopFile())
		if err != nil {
			logger.Error(err)
		} else {
			aliasArgs := []string{
				"--force",
				"-n", appInstalled.GetRegisterName(),
				"-c", execCmdStr,
			}
			libs.RunCoreUtilsCmd("alias-manager-cu", false, aliasArgs...)
		}
	}
	logger.Error(file.DeleteFile(packageFileTemp))
}

func main() {
	cobralib.Run()
}
