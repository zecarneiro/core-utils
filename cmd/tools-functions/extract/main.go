package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	fileArg     string
	destination string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "extract <file>",
		Short: "Extract an archive file",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().StringVarP(&destination, "destination", "d", "", "Destination directory")
	cobralib.WithRunArgsStr(process)
}

func process(filePath string) {
	if str.IsEmpty(destination) {
		cwd, err := file.GetCurrentDir()
		logic.ProcessError(err)
		destination = cwd
	}
	if !file.IsFile(filePath) {
		logic.ProcessError(fmt.Errorf("given file not found: %s", filePath))
	}
	ext := file.FileExtension(filePath)
	cmd := getExtractCmd(filePath, ext)
	if str.IsEmpty(cmd) {
		logic.ProcessError(fmt.Errorf("don't know how to extract '%s'...", filePath))
	}
	execCmd := models.Command{
		Cmd:      cmd,
		Verbose:  true,
		UseShell: true,
	}
	logic.ProcessError(exe.ExecRealTime(execCmd))
}

func getExtractCmd(filePath, ext string) string {
	// Windows
	if platform.IsWindows() {
		switch ext {
		case "zip":
			return fmt.Sprintf("Expand-Archive -LiteralPath \"%s\" -DestinationPath \"%s\"", filePath, destination)
		default:
			return ""
		}
	}
	// Others Platforms
	switch ext {
	case ".zip":
		return fmt.Sprintf("unzip \"%s\" -d \"%s\"", filePath, destination)
	case ".tar.bz2", ".tbz2":
		return fmt.Sprintf("tar xvjf \"%s\"", filePath)
	case ".tar.gz", ".tgz":
		return fmt.Sprintf("tar xvzf \"%s\"", filePath)
	case ".bz2":
		return fmt.Sprintf("bunzip2 \"%s\"", filePath)
	case ".rar":
		return fmt.Sprintf("rar x \"%s\"", filePath)
	case ".gz":
		return fmt.Sprintf("gunzip \"%s\"", filePath)
	case ".tar":
		return fmt.Sprintf("tar xvf \"%s\"", filePath)
	case ".Z":
		return fmt.Sprintf("uncompress \"%s\"", filePath)
	case ".7z":
		return fmt.Sprintf("7z x \"%s\"", filePath)
	}
	return ""
}

func main() {
	cobralib.Run()
}
