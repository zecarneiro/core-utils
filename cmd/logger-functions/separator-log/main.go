package main

import (
	"fmt"
	"golangutils/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	length  int
	rootCmd *cobra.Command

	defaultLen = 50
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "separator-log",
		Short: "Print a separator by length",
	}
	rootCmd.Flags().IntVarP(&length, "length", "l", defaultLen, fmt.Sprintf("Length of separator. Default is %d", defaultLen))
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		logger.WithSeparatorLength(length)
		logger.Separator()
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
