package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "lhidenf",
		Short: "List hidden files",
	}
	cobralib.WithWorkingDirDefault()
	cobralib.WithRun(process)
}

func process() {
	workingDir := cobralib.GetWorkingDir()
	filesInfo, err := file.ReadDirRecursive(workingDir)
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
