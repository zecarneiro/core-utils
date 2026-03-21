package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var isRecursive bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "listdir",
		Short: "List directories in the current directory",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&isRecursive, "recursive", "r", false, "List directories recursively")
	cobralib.WithWorkingDirDefault()
	cobralib.WithRun(process)
}

func process() {
	var filesInfo models.FileInfo
	var err error
	workingDir := cobralib.GetWorkingDir()
	if isRecursive {
		filesInfo, err = file.ReadDirRecursive(workingDir)
	} else {
		filesInfo, err = file.ReadDir(workingDir)
	}
	logic.ProcessError(err)
	for _, directory := range filesInfo.Directories {
		fmt.Println(directory)
	}
}

func main() {
	cobralib.Run()
}
