package main

import (
	"fmt"
	"golangutils/pkg/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	msgParts []string
	rootCmd  *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "title-log",
		Short: "Print a title from message parts",
	}
	rootCmd.Flags().StringSliceVarP(&msgParts, "message", "m", []string{}, "Message parts. Can be repeated: -m part1 -m part2, or comma-separated: -m \"part1,part2\"")
	_ = rootCmd.MarkFlagRequired("message")
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		logger.Title(strings.Join(msgParts, " "))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
