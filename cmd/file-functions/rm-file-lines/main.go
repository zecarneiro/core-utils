package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	filepath string
	match    string
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
	cobralib.WithRun(process)
}

func filterFileData(fileLine string, err error) {
	if !strings.Contains(fileLine, match) {
		data += logic.Ternary(data == "", fileLine, common.Eol()+fileLine)
	}
	logic.ProcessError(err)
}

func process() {
	if match == "" {
		logic.ProcessError(fmt.Errorf("invalid given match"))
	}
	file.ReadFileLineByLine(filepath, filterFileData)
	logic.ProcessError(file.WriteFile(filepath, data, false, false))
}

func main() {
	cobralib.Run()
}
