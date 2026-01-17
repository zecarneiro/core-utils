package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"os"

	"github.com/spf13/cobra"
)

var (
	filepath          string
	match             string
	isCaseInsensitive bool
	rootCmd           *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "file-contain",
		Short: "Check if file contain given match",
	}
	rootCmd.Flags().StringVarP(&filepath, "file", "f", "", "File to write content")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringVarP(&match, "match", "m", "", "Match content to verify")
	rootCmd.MarkFlagRequired("match")
	rootCmd.Flags().BoolVarP(&isCaseInsensitive, "case-insensitive", "i", false, "Enable search with case insensitive")
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if match == "" {
			generic.ProcessError(fmt.Errorf("invalid given match"))
			os.Exit(1)
		}
		data, err := file.ReadFile(filepath)
		generic.ProcessError(err)
		fmt.Println(common.StringContains(data, match, isCaseInsensitive))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
