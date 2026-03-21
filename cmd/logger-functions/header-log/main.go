package main

import (
	"fmt"
	"golangutils/pkg/logger"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var length int

func init() { setupCommand() }

func setupCommand() {
	defaultLen := 50
	cobralib.CobraCmd = &cobra.Command{
		Use:   "header-log [message]",
		Short: "Print a header from message",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().IntVarP(&length, "length", "l", defaultLen, fmt.Sprintf("Length of header. Default is %d", defaultLen))
	cobralib.WithRunArgsStr(process)
}

func process(msg string) {
	logger.WithHeaderLength(length)
	logger.Header(msg)
}

func main() {
	cobralib.Run()
}
