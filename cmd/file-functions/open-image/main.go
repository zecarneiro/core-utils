package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "open-image [filepath]",
		Short: "Open image",
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
		cmdInfo = models.Command{Cmd: "Start-Process", Args: []string{filepath}, UseShell: true}
	} else if platform.IsLinux() {
		cmdInfo = models.Command{Cmd: "xdg-open", Args: []string{filepath}, UseShell: false}
	} else if platform.IsDarwin() {
		cmdInfo = models.Command{Cmd: "open", Args: []string{filepath}, UseShell: true}
	}
	cmdInfo.Verbose = false
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
}

func main() {
	cobralib.Run()
}
