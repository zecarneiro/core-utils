package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"
	"os"

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
	if match == "" {
		logic.ProcessError(fmt.Errorf("invalid given match"))
		os.Exit(1)
	}
	data, err := file.ReadFile(filepath)
	logic.ProcessError(err)
	fmt.Println(str.StringContains(data, match, isCaseInsensitive))
}

func main() {
	cobralib.Run()
}
