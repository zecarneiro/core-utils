package main

import (
	"golangutils/pkg/logger"
	"golangutils/pkg/netc"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var verbose bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "has-internet",
		Short: "Check if there is internet connection",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show all error output during check the internet")
	cobralib.WithRun(process)
}

func process() {
	status, err := netc.HasInternet()
	if verbose {
		logger.Error(err)
	}
	logger.Log(status)
}

func main() {
	cobralib.Run()
}
