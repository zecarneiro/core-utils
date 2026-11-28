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
		Short: "Install AppImage package using Gear Lever.",
	}
	cobralib.WithRunArgsStr(process)
}

func installAllCmds() {
	logger.Info("Creating terminal comand...")
	libs.RunCoreUtilsCmdWithOutput("appimage-install-cmds", false)
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
	logger.Info(fmt.Sprintf("Deleting copied file: %s", packageFileTemp))
	logger.Error(file.DeleteFile(packageFileTemp))
	installAllCmds()
	logger.Ok("Installation is done.")
}

func main() {
	cobralib.Run()
}
