package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	filepath string
	match    string
	isRegex  bool
	data     = ""
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "rm-file-lines",
		Short: "Delete lines of file that match with given argument",
	}
	cobralib.CobraCmd.Flags().StringVarP(&filepath, "file", "f", "", "File to write content")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("file"))
	cobralib.CobraCmd.Flags().StringVarP(&match, "match", "m", "", "Match content in lines to delete")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("match"))
	cobralib.CobraCmd.Flags().BoolVarP(&isRegex, "is-regex", "r", false, "Apply Match content as regex")
	cobralib.WithRun(process)
}

func process() {
	if str.IsEmpty(match) {
		logic.ProcessError(fmt.Errorf("invalid given match"))
	}
	logic.ProcessError(file.DeleteFileLines(filepath, match, isRegex))
}

func main() {
	cobralib.Run()
}
