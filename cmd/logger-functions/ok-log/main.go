package main

import (
	"fmt"
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
		Use:   "ok-log [message]",
		Short: "Print an OK log message",
		Args:  cobra.MinimumNArgs(1),
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		data := strings.Join(args, " ")
		logger.Ok(data)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
