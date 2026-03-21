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
	filepath          string
	match             string
	isCaseInsensitive bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "file-contain",
		Short: "Check if file contain given match",
	}
	cobralib.CobraCmd.Flags().StringVarP(&filepath, "file", "f", "", "File to write content")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("file"))
	cobralib.CobraCmd.Flags().StringVarP(&match, "match", "m", "", "Match content to verify")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("match"))
	cobralib.CobraCmd.Flags().BoolVarP(&isCaseInsensitive, "case-insensitive", "i", false, "Enable search with case insensitive")
	cobralib.WithRun(process)
}

func process() {
	if str.IsEmpty(match) {
		logic.ProcessError(fmt.Errorf("invalid given match"))
	}
	status, err := file.FileTextContains(filepath, match, isCaseInsensitive)
	logic.ProcessError(err)
	fmt.Println(status)
}

func main() {
	cobralib.Run()
}
