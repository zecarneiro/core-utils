package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
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
		Use:   "file-encoding [file]",
		Short: "Get encoding of file",
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
		name, err := file.GetFileEncoding(data)
		generic.ProcessError(err)
		fmt.Println(name)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
