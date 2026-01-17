package main

import (
	"fmt"
	"golangutils/pkg/entity"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"os"

	"github.com/spf13/cobra"
)

var (
	isRecursive bool
	rootCmd     *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "count-dir",
		Short: "Count directories in the current directory",
	}
	rootCmd.Flags().BoolVarP(&isRecursive, "recursive", "r", false, "Count directories recursively")
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		var filesInfo entity.FileInfo
		currentDir, err := file.GetCurrentDir()
		generic.ProcessError(err)
		if isRecursive {
			filesInfo, err = file.ReadDirRecursive(currentDir)
		} else {
			filesInfo, err = file.ReadDir(currentDir)
		}
		generic.ProcessError(err)
		fmt.Println(len(filesInfo.Directories))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
