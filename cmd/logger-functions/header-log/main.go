package main

import (
	"fmt"
	"golangutils/pkg/logger"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	msgParts   []string
	length     int
	defaultLen = 50
	rootCmd    *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "header-log",
		Short: "Print a header from message parts",
	}
	rootCmd.Flags().StringSliceVarP(&msgParts, "message", "m", []string{}, "Message parts. Can be repeated: -m part1 -m part2, or comma-separated: -m \"part1,part2\"")
	_ = rootCmd.MarkFlagRequired("message")

	rootCmd.Flags().IntVarP(&length, "length", "l", defaultLen, fmt.Sprintf("Length of header. Default is %d", defaultLen))
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		logger.WithHeaderLength(length)
		logger.Header(strings.Join(msgParts, " "))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
