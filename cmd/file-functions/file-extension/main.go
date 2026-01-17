package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "file-extension [file]",
		Short: "Get extension of file",
		Args:  cobra.MinimumNArgs(1),
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		data := strings.Join(args, " ")
		data = common.Ternary(data == "." || data == "..", "", data)
		if !file.IsFile(data) {
			logger.Error(fmt.Errorf("%s is not a file or file not found", data))
			os.Exit(1)
		}
		fmt.Println(file.FileExtension(data))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
