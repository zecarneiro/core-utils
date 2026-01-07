package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/system"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	if !system.IsAdmin() {
		logic.ProcessError(fmt.Errorf("Please, need root access to continue!"))
	}
	cobralib.CobraCmd = &cobra.Command{
		Use:   "snap-uninstall [app]",
		Short: "Uninstall snap app",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func deleteDir(dir string) {
	if file.IsDir(dir) {
		logger.Info("Deleting directory: " + dir)
		logic.ProcessError(file.DeleteDirectory(dir))
	}
}

func process(app string) {
	cmdInfo := models.Command{
		Verbose:  true,
		UseShell: true,
	}
	configDir := file.ResolvePath(system.HomeDir(), "snap", app)
	configSystemDir := file.ResolvePath("/snap", app)
	cmdInfo.Cmd = fmt.Sprintf("sudo snap remove --purge %s", app)
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
	cmdInfo.Cmd = "snap saved"
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
	fmt.Print("Insert the number on the line of App(ENTER TO SKIP): ")
	var response string
	fmt.Scanln(&response)
	response = strings.Trim(response, " ")
	if len(response) > 0 {
		cmdInfo.Cmd = fmt.Sprintf("sudo snap forget %s", response)
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	}
	deleteDir(configDir)
	deleteDir(configSystemDir)
}

func main() {
	cobralib.Run()
}
