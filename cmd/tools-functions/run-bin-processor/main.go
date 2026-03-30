package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	fileArgs     string
	forceArg     bool
	acceptExtMsg = "Only accept file with ext: .exe|.msi|.msixbundle"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "run-bin-processor <file_or_directory>",
		Short: fmt.Sprintf("Run a binary file or all binaries in a directory (Windows only - %s)", acceptExtMsg),
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().StringVarP(&fileArgs, "args", "a", "", "File args in one string. If pass directory, each bin file will receive the same args")
	cobralib.CobraCmd.Flags().BoolVarP(&forceArg, "force", "f", false, "Enable force flag for AppxPackage")
	cobralib.WithRunArgsStr(process)
}

func process(fileOrDir string) {
	fileOrDir = file.ResolvePath(fileOrDir)
	if str.IsEmpty(fileOrDir) {
		logic.ProcessError(errors.New("Please give a valid file or directory argument"))
	}
	if file.IsFile(fileOrDir) {
		runBin(fileOrDir)
	} else if file.IsDir(fileOrDir) {
		runDir(fileOrDir)
	} else {
		logic.ProcessError(fmt.Errorf("invalid given bin file or directory. Only work for .exe and .msi files"))
	}
}

func runBin(binary string) {
	var execCmd models.Command
	ext := strings.ToLower(file.FileExtension(binary))

	switch ext {
	case "exe", "msi":
		cmd := fmt.Sprintf("Start-Process '%s' -Wait", binary)
		if !str.IsEmpty(fileArgs) {
			cmd = fmt.Sprintf("%s -ArgumentList %s", cmd, fileArgs)
		}
		execCmd = models.Command{
			Cmd:      cmd,
			UseShell: true,
			Verbose:  true,
		}
	case "msixbundle", "appx":
		execCmd = models.Command{
			Cmd:      fmt.Sprintf("Add-AppxPackage -Path \"%s\" %s", binary, logic.Ternary(forceArg, "-ForceApplicationShutdown", "")),
			UseShell: true,
			Verbose:  true,
		}
	default:
		logger.Info(acceptExtMsg)
		logic.ProcessError(fmt.Errorf("can not run this bin file: %s", binary))
	}
	logic.ProcessError(exe.ExecRealTime(execCmd))
}

func runDir(directory string) {
	files, err := file.ReadDir(directory)
	logic.ProcessError(err)
	for _, fileinfo := range files.Files {
		fileinfo = file.JoinPath(directory, fileinfo)
		runBin(fileinfo)
	}
	logger.Info(fmt.Sprintf("Execution of all files on '%s' it's done.", directory))
}

func main() {
	cobralib.Run()
}
