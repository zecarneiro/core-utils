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
		Use:   "findfile [filename]",
		Short: "Find a file by name",
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
		for _, filepath := range filesInfo.Files {
			filename := file.Basename(filepath)
			if strings.Contains(filename, data) {
				fmt.Println(filepath)
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
