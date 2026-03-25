package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "change-user-default-dir",
		Short: "",
	}
	cobralib.WithRun(process)
}

func process() {
	if console.Confirm("Insert all User Dirs?", true) {
		args := []string{}
		scriptExtension := ".sh"
		if platform.IsWindows() {
			err := errors.New("Can not get shell profile dir")
			scriptExtension = ".ps1"
			currentScriptFile := libs.RunCoreUtilsCmdWithOutput("SHELL_PROFILE_SCRIPT", false)
			if str.IsEmpty(currentScriptFile) {
				logic.ProcessError(err)
			}
			currentScriptDir := file.Dirname(currentScriptFile)
			if str.IsEmpty(currentScriptDir) || !file.IsDir(currentScriptDir) {
				logic.ProcessError(err)
			}
			args = append(args, currentScriptDir)
		}
		scriptFile := libs.GetScriptCmdPathByName(fmt.Sprintf("change-user-default-dir%s", scriptExtension), "system")
		cmd := models.Command{
			Cmd:      scriptFile,
			Args:     args,
			UseShell: true,
		}
		logic.ProcessError(exe.ExecRealTime(cmd))
	}
}

func main() {
	cobralib.Run()
}
