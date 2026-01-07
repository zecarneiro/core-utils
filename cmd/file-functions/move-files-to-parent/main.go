package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/slice"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "move-files-to-parent",
		Short: "Move all file on given directory to main directory",
	}
	cobralib.WithWorkingDirDefault()
	cobralib.WithRun(process)
}

func process() {
	workingDir := cobralib.GetWorkingDir()
	filesInfo, err := file.ReadDirRecursive(workingDir)
	logic.ProcessError(err)
	fileList := slice.FilterArray(filesInfo.Files, func(filepath string) bool {
		dirname := file.Dirname(filepath)
		return dirname != workingDir
	})
	fileListSize := len(fileList)
	for index, filepath := range fileList {
		dirname := file.Dirname(filepath)
		if dirname != workingDir {
			filename := file.Basename(filepath)
			destFile := file.ResolvePath(workingDir, filename)
			canMove := true
			if file.IsFile(destFile) && !console.Confirm(fmt.Sprintf("Already exist this file: %s. Continue?", destFile), true) {
				canMove = false
			}
			if canMove {
				logger.Info(fmt.Sprintf("Moving (%d of %d): %s ...", index+1, fileListSize, filepath))
				err := file.Move(filepath, workingDir)
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}
}

func main() {
	cobralib.Run()
}
