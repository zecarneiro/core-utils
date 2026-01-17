package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"golangutils/pkg/logger"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "rm-empty-dir [directory_path]",
		Short: "Delete empty directories",
	}
}

func showErrorOnReadDir(dir string, err error, isDelete bool) {
	if err != nil {
		msg := common.Ternary(isDelete, "Get error on delete directory:", "Get error on read directory:")
		logger.Error(msg + " " + dir)
		logger.Error(err)
	}
}

func deleteFilepathInCurrentDir(fileList []string, currentDir string) []string {
	newFileList := []string{}
	for _, filepath := range fileList {
		dirname := file.Dirname(filepath)
		if dirname != currentDir {
			newFileList = append(newFileList, filepath)
		}
	}
	return newFileList
}

func process(data string) {
	filesInfo, err := file.ReadDirRecursive(data)
	if err != nil {
		showErrorOnReadDir(data, err, false)
	} else {
		filesInfo.Directories = common.FilterArray(filesInfo.Directories, func(val string) bool {
			return val != data
		})
		slices.Reverse(filesInfo.Directories)
		if len(filesInfo.Files) == 0 && len(filesInfo.Directories) == 0 {
			err = file.DeleteDirectory(data)
			if err != nil {
				showErrorOnReadDir(data, err, true)
			} else {
				logger.Ok("Directory was deleted: " + data)
			}
		} else {
			for _, dirInfo := range filesInfo.Directories {
				if dirInfo != data {
					process(dirInfo)
				}
			}
		}
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		msgErr := "Invalid given directory path"
		data := strings.Join(args, " ")
		if data == ".." {
			logger.Error(msgErr)
			os.Exit(1)
		}
		if data == "." || len(data) == 0 {
			currentDir, err := file.GetCurrentDir()
			generic.ProcessError(err)
			data = currentDir
		}
		data = file.ResolvePath(data)
		if !file.IsDir(data) {
			logger.Error(msgErr)
			os.Exit(1)
		}
		process(data)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
