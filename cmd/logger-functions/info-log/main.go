package main

import (
	"golangutils/pkg/logger"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "info-log [message]",
		Short: "Print an info log message",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(msg string) {
	logger.Info(msg)
}

func main() {
	cobralib.Run()
}
