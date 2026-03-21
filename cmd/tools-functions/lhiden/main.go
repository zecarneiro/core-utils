package main

import (
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var isRecursive bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "lhidend",
		Short: "List hidden directories and files",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&isRecursive, "recursive", "r", false, "List hidden directories recursively")
	cobralib.WithWorkingDirDefault()
	cobralib.WithRun(process)
}

func process() {
	recursiveFlag := logic.Ternary(isRecursive, "-r", "")
	logger.Header("List hidden files")
	libs.RunCoreUtilsCmd("lhidenf", false, recursiveFlag)
	logger.Header("List hidden directories")
	libs.RunCoreUtilsCmd("lhidend", false, recursiveFlag)
}

func main() {
	cobralib.Run()
}
