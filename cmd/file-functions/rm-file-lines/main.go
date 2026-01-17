package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"golangutils/pkg/system"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	filepath string
	match    string
	rootCmd  *cobra.Command

	data = ""
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "rm-file-lines",
		Short: "Delete lines of file that match with given argument",
	}
	rootCmd.Flags().StringVarP(&filepath, "file", "f", "", "File to write content")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringVarP(&match, "match", "m", "", "Match content in lines to delete")
	rootCmd.MarkFlagRequired("match")
}

func filterFileData(fileLine string, err error) {
	if !strings.Contains(fileLine, match) {
		data += common.Ternary(data == "", fileLine, system.Eol()+fileLine)
	}
	generic.ProcessError(err)
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if match == "" {
			generic.ProcessError(fmt.Errorf("invalid given match"))
			os.Exit(1)
		}
		file.ReadFileLineByLine(filepath, filterFileData)
		generic.ProcessError(file.WriteFile(filepath, data, false, false))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
