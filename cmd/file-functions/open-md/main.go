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

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "open-md <filepath>",
		Short: "Open and view markdown file",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(filepath string) {
	var cmdInfo models.Command
	if filepath == "." || filepath == ".." {
		filepath = ""
	}
	if !file.IsFile(filepath) {
		logic.ProcessError(fmt.Errorf("Invalid given file: %s", filepath))
	}
	if platform.IsWindows() {
		cmdPath, err := console.Which("mdview")
		logic.ProcessError(err)
		libs.RunCoreUtilsCmd("run-bin-processor", false, cmdPath, "-a", filepath)
	} else if platform.IsLinux() {
		cmdInfo = models.Command{Cmd: "markdown_viewer", Args: []string{filepath}, UseShell: false}
	} else {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
	cmdInfo.Verbose = false
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
}

func main() {
	cobralib.Run()
}
