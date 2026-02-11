package main

import (
	"golangutils/pkg/logger"
	"golangutils/pkg/system"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "ver",
		Short: "SO Info",
	}
	cobralib.WithRun(process)
}

func process() {
	logger.Log(system.OSName())
}

func main() {
	cobralib.Run()
}
