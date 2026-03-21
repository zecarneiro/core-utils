package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	pathArg string
	typeArg string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "unix2win-path-conv",
		Short: "Convert paths between Unix and Windows format (Windows only)",
	}
	cobralib.CobraCmd.Flags().StringVarP(&pathArg, "path", "p", "", "Path to convert (required)")
	cobralib.CobraCmd.Flags().StringVarP(&typeArg, "type", "t", "", "Conversion type: unix2win or win2unix (required)")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("path"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("type"))
	cobralib.WithRun(process)
}

func process() {
	var cmd string
	if str.IsEmpty(pathArg) {
		logic.ProcessError(errors.New("insert a valid path arg"))
	}
	switch typeArg {
	case "unix2win":
		cmd = fmt.Sprintf("cygpath -w \"%s\"", pathArg)
	case "win2unix":
		cmd = fmt.Sprintf("cygpath -u \"%s\"", pathArg)
	default:
		logic.ProcessError(errors.New("insert a valid type arg"))
	}
	execCmd := models.Command{
		Cmd:      cmd,
		UseShell: true,
	}
	logic.ProcessError(exe.ExecRealTime(execCmd))
}

func main() {
	cobralib.Run()
}
