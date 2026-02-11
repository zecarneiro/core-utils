package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	extensionArg string
	typeArg      string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "dos2unix-recursive",
		Short: "Convert files recursively between DOS and UNIX line endings",
	}
	cobralib.CobraCmd.Flags().StringVarP(&extensionArg, "extension", "e", "", "File extension to convert (required)")
	cobralib.CobraCmd.Flags().StringVarP(&typeArg, "type", "t", "", "Conversion type: dos2unix or unix2dos (required)")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("extension"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("type"))
	cobralib.WithRun(process)
}

func process() {
	if extensionArg == "" {
		return
	}
	var cmd string
	if platform.IsWindows() {
		cmd = fmt.Sprintf("Get-ChildItem -Recurse -File -Filter \"*.%s\" | ForEach-Object ", extensionArg)
		switch typeArg {
		case "dos2unix":
			cmd = cmd + "{ dos2unix $_.FullName }"
		case "unix2dos":
			cmd = cmd + "{ unix2dos $_.FullName }"
		default:
			logic.ProcessError(errors.New("Invalid given of type arg"))
		}
	} else {
		cmd = fmt.Sprintf("find . -type f -name \"*.%s\" -print0 | xargs -0 ", extensionArg)
		switch typeArg {
		case "dos2unix":
			cmd = cmd + "dos2unix"
		case "unix2dos":
			cmd = cmd + "unix2dos"
		default:
			logic.ProcessError(errors.New("Invalid given of type arg"))
		}
	}
	execCmd := models.Command{
		Cmd:      cmd,
		Verbose:  true,
		UseShell: true,
	}
	err := exe.ExecRealTime(execCmd)
	if err != nil {
		logger.Error(err)
	}
}

func main() {
	cobralib.Run()
}
