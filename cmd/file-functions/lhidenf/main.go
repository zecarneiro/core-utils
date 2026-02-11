package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var isRecursive bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "lhidenf",
		Short: "List hidden files",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&isRecursive, "recursive", "r", false, "List hidden files recursively")
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
	for _, filepath := range filesInfo.Files {
		isHidden, err := file.IsHidden(filepath)
		if err != nil {
			logger.Error(err)
		} else if isHidden {
			fmt.Println(filepath)
		}
	}
}

func main() {
	cobralib.Run()
}
