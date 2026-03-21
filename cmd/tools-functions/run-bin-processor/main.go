package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/console"
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

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "run-bin-processor <file_or_directory>",
		Short: "Run a binary file or all binaries in a directory (Windows only)",
		Args:  cobra.MinimumNArgs(1),
	}
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
		logic.ProcessError(fmt.Errorf("invalid given bin file or directory"))
	}
	console.Pause()
}

func runBin(binary string) {
	var execCmd models.Command
	ext := strings.ToLower(file.FileExtension(binary))

	switch ext {
	case "exe", "msi":
		execCmd = models.Command{
			Cmd:      fmt.Sprintf("Start-Process '%s' -Wait", binary),
			UseShell: true,
			Verbose:  true,
		}
	case "msixbundle":
		execCmd = models.Command{
			Cmd:      fmt.Sprintf("Add-AppxPackage -Path \"%s\"", binary),
			UseShell: true,
			Verbose:  true,
		}
	default:
		logger.Info("Only accept file with ext: .exe|.msi|.msixbundle")
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
