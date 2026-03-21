package main

import (
	"golangutils/pkg/logger"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "debug-log [message]",
		Short: "Print a debug log message",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(msg string) {
	logger.Debug(msg)
}

func main() {
	cobralib.Run()
}
