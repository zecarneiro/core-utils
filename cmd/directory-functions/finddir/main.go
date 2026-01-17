package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "finddir [directory_name]",
		Short: "Find a directory by name",
		Args:  cobra.MinimumNArgs(1),
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		currentDir, err := file.GetCurrentDir()
		generic.ProcessError(err)
		data := strings.Join(args, " ")
		filesInfo, err := file.ReadDirRecursive(currentDir)
		generic.ProcessError(err)
		for _, directory := range filesInfo.Directories {
			dirname := file.Basename(directory)
			if strings.Contains(dirname, data) {
				fmt.Println(directory)
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
