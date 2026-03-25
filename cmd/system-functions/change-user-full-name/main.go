package main

import (
	"errors"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "change-user-full-name",
		Short: "",
	}
	cobralib.WithRun(process)
}

func process() {
	if !platform.IsWindows() {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
	scriptFile := libs.GetScriptCmdPathByName("change-user-full-name.ps1", "system")
	cmd := models.Command{
		Cmd:      scriptFile,
		UseShell: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
