package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	messageArg  []string
	isNoDefault bool
)

func init() {
	console.EnableFeatures()
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "confirm",
		Short: "Show confirm message and request a response from user input",
	}
	cobralib.CobraCmd.Flags().StringArrayVarP(&messageArg, "message", "m", []string{}, "Question to user (required)")
	cobralib.CobraCmd.Flags().BoolVarP(&isNoDefault, "is-no-default", "n", false, "If set, the default response will be 'false'")
	cobralib.WithRun(process)
}

func process() {
	messageStr := slice.ArrayToString(messageArg)
	if str.IsEmpty(messageStr) {
		logic.ProcessError(fmt.Errorf("Message can not be empty"))
	}
	response := console.Confirm(messageStr, isNoDefault)
	logger.Log(response)
}

func main() {
	cobralib.Run()
}
