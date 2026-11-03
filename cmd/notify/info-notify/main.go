package main

import (
	"golangutils/pkg/logic"
	"golangutils/pkg/ui"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	titleArg   string
	messageArg string
	iconArg    string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "info-notify",
		Short: "Print an info notify message",
	}
	cobralib.CobraCmd.Flags().StringVarP(&titleArg, "title", "t", "", "Title to show")
	cobralib.CobraCmd.Flags().StringVarP(&messageArg, "message", "m", "", "Message to show")
	cobralib.CobraCmd.Flags().StringVarP(&iconArg, "icon", "i", "", "Icon to show")
	cobralib.WithRun(process)
}

func process() {
	logic.ProcessError(ui.InfoNofity(titleArg, messageArg, iconArg))
}

func main() {
	cobralib.Run()
}
