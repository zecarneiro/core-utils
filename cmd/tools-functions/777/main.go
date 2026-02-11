package main

import (
	"errors"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() {
	console.EnableFeatures()
	if !system.IsAdmin() && platform.IsWindows() {
		logic.ProcessError(errors.New(system.NeedAdminAccessMsg))
	}
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "777 <path>",
		Short: "Set full permission for a file or directory",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(path string) {
	logic.ProcessError(exe.Chmod777(path, true))
}

func main() {
	cobralib.Run()
}
