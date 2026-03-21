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
		Use:   "separator-log",
		Short: "Print a separator by length",
	}
	cobralib.CobraCmd.Flags().IntVarP(&length, "length", "l", defaultLen, fmt.Sprintf("Length of separator. Default is %d", defaultLen))
	cobralib.WithRun(process)
}

func process() {
	logger.WithSeparatorLength(length)
	logger.Separator()
}

func main() {
	cobralib.Run()
}
