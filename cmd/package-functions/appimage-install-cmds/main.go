package main

import (
	"fmt"
	"golangutils/pkg/logger"
	"golangutils/pkg/system"
	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "appimage-install-cmds",
		Short: "Install AppImage packages commands.",
	}
	cobralib.WithRun(process)
}

func process() {
	for _, appInstalled := range libs.GetGearLeverListOfInstalledApps().Installed {
		if !appInstalled.IsNotRegisted {
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
		} else {
			logger.Warn(fmt.Sprintf(`Not processed: %s`, appInstalled.Name))
		}
	}
	logger.Ok("Installation is done.")
}

func main() {
	cobralib.Run()
}
