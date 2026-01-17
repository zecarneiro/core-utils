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
		Use:   "lhiden-dir",
		Short: "List hidden directories",
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		currentDir, err := file.GetCurrentDir()
		generic.ProcessError(err)
		filesInfo, err := file.ReadDirRecursive(currentDir)
		generic.ProcessError(err)
		for _, directory := range filesInfo.Directories {
			isHidden, err := file.IsHidden(directory)
			if err != nil {
				logger.Error(err)
			} else if isHidden {
				fmt.Println(directory)
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
