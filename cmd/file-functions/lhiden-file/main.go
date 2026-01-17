package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"golangutils/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "lhiden-file",
		Short: "List hidden files",
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		currentDir, err := file.GetCurrentDir()
		generic.ProcessError(err)
		filesInfo, err := file.ReadDirRecursive(currentDir)
		generic.ProcessError(err)
		for _, filepath := range filesInfo.Files {
			isHidden, err := file.IsHidden(filepath)
			if err != nil {
				logger.Error(err)
			} else if isHidden {
				fmt.Println(filepath)
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
