package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"slices"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "rm-empty-dir",
		Short: "Delete empty directories",
	}
	cobralib.WithWorkingDirDefault()
	cobralib.WithRun(process)
}

func showErrorOnReadDir(dir string, err error, isDelete bool) {
	if err != nil {
		msg := logic.Ternary(isDelete, "Get error on delete directory:", "Get error on read directory:")
		logger.Error(msg + " " + dir)
		logger.Error(err)
	}
}

func deleteAll(path string) {
	path = file.ResolvePath(path)
	if !file.IsDir(path) {
		logic.ProcessError(fmt.Errorf("Invalid given directory path"))
	}
	filesInfo, err := file.ReadDirRecursive(path)
	showErrorOnReadDir(path, err, false)
	slices.Reverse(filesInfo.Directories)
	if len(filesInfo.Files) == 0 && len(filesInfo.Directories) == 0 {
		err = file.DeleteDirectory(path)
		showErrorOnReadDir(path, err, true)
		logger.Ok("Directory was deleted: " + path)
	} else {
		for _, dirInfo := range filesInfo.Directories {
			if dirInfo != path {
				deleteAll(dirInfo)
			}
		}
	}
}

func process() {
	deleteAll(cobralib.GetWorkingDir())
}

func main() {
	cobralib.Run()
}
