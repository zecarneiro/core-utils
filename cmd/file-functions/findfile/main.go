package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "findfile [filename]",
		Short: "Find a file by name",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithWorkingDirUsage("Working directory to search")
	cobralib.WithRunArgsStr(process)
}

func process(search string) {
	workingDir := cobralib.GetWorkingDir()
	filesInfo, err := file.ReadDirRecursive(workingDir)
	logic.ProcessError(err)
	for _, filepath := range filesInfo.Files {
		filename := file.Basename(filepath)
		if strings.Contains(filename, search) {
			fmt.Println(filepath)
		}
	}
}

func main() {
	cobralib.Run()
}
