package main

import (
	"fmt"
	"golangutils/pkg/logic"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	statusArg bool
	valueArg  int
)

var config *libs.Config

func init() {
	config = libs.NewConfig()
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "prompt-style",
		Short: "Get or set prompt style",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&statusArg, "status", "s", false, "Return current value of prompt style")
	cobralib.CobraCmd.Flags().IntVarP(&valueArg, "value", "v", config.PromptStyle, "Change prompt style (1-4)")
	cobralib.WithRun(process)
}

func process() {
	if statusArg {
		fmt.Println(config.PromptStyle)
	} else if config.PromptStyle != valueArg {
		config.PromptStyle = valueArg
		config.Write()
	} else {
		logic.ProcessError(cobralib.CobraCmd.Help())
	}
}

func main() {
	cobralib.Run()
}
