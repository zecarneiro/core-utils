package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "file-exists [filepath]",
		Short: "Check if a file exists",
		Args:  cobra.MinimumNArgs(1),
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		data := strings.Join(args, " ")
		data = common.Ternary(data == "." || data == "..", "", data)
		fmt.Println(file.IsFile(data))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
