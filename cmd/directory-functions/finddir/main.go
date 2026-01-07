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
		Use:   "finddir [directory_name]",
		Short: "Find a directory by name",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithWorkingDirUsage("Working directory to search")
	cobralib.WithRunArgsStr(process)
}

func process(search string) {
	workingDir := cobralib.GetWorkingDir()
	filesInfo, err := file.ReadDirRecursive(workingDir)
	logic.ProcessError(err)
	for _, directory := range filesInfo.Directories {
		dirname := file.Basename(directory)
		if strings.Contains(dirname, search) {
			fmt.Println(directory)
		}
	}
}

func main() {
	cobralib.Run()
}
