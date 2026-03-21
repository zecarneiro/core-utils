package main

import (
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	scriptPathArg string
	scriptArgsArg []string
	shellTypeArg  string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "run-shell-script",
		Short: "Run shell script",
	}
	cobralib.CobraCmd.Flags().StringVarP(&scriptPathArg, "file", "f", "", "Script file path (required)")
	cobralib.CobraCmd.Flags().StringArrayVarP(&scriptArgsArg, "args", "a", []string{}, "Args for script")
	cobralib.CobraCmd.Flags().StringVarP(&shellTypeArg, "shell", "s", "", "Define type of script is and define wich shell need to use (required)")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("file"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("shell"))
	cobralib.WithRun(process)
}

func process() {
	shellType := enums.GetShellTypeFromValue(shellTypeArg)
	logic.ProcessError(exe.RunScriptShell(scriptPathArg, shellType, false, scriptArgsArg...))
}

func main() {
	cobralib.Run()
}
