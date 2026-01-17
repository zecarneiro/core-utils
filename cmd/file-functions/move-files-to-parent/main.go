package main

import (
	"fmt"
	"golangutils/pkg/common"
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
		Use:   "move-files-to-parent",
		Short: "Move all file on given directory to main directory",
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		currentDir, err := file.GetCurrentDir()
		generic.ProcessError(err)
		filesInfo, err := file.ReadDirRecursive(currentDir)
		generic.ProcessError(err)
		fileList := common.FilterArray(filesInfo.Files, func(filepath string) bool {
			dirname := file.Dirname(filepath)
			return dirname != currentDir
		})
		fileListSize := len(fileList)
		for index, filepath := range fileList {
			dirname := file.Dirname(filepath)
			if dirname != currentDir {
				filename := file.Basename(filepath)
				destFile := file.ResolvePath(currentDir, filename)
				canMove := true
				if file.IsFile(destFile) && !generic.Confirm(fmt.Sprintf("Already exist this file: %s. Continue?", destFile), true) {
					canMove = false
				}
				if canMove {
					logger.Info(fmt.Sprintf("Moving (%d of %d): %s ...", index+1, fileListSize, filepath))
					err := file.Move(filepath, currentDir)
					if err != nil {
						logger.Error(err)
					}
				}
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
