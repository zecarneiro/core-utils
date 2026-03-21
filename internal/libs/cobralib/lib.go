package cobralib

import (
	"golangutils/pkg/logic"

	"github.com/spf13/cobra"
)

var (
	CobraCmd   *cobra.Command
	workingDir string
	cobraArgs  []string

	enableFlagParsing = false
)

var (
	FuncRun     func(cmd *cobra.Command, args []string)
	FuncExecute func() error
)

func Run() {
	if FuncExecute == nil {
		setFuncs()
	}
	if err := FuncExecute(); err != nil {
		logic.ProcessError(err)
	}
}
