package cobralib

import (
	"golangutils/pkg/logic"
	"golangutils/pkg/slice"

	"main/internal/libs"

	"github.com/spf13/cobra"
)

func setFuncs() {
	CobraCmd.Run = FuncRun
	FuncExecute = CobraCmd.Execute
}

func withCobraArgs(args []string) {
	cobraArgs = args
}

func WithWorkingDir(name string, short string, defaultVal string, usage string) {
	CobraCmd.Flags().StringVarP(&workingDir, name, short, defaultVal, usage)
}

func WithWorkingDirUsage(usage string) {
	WithWorkingDir("working-dir", "w", "", logic.Ternary(len(usage) > 0, usage, "Working directory"))
}

func WithWorkingDirDefault() {
	WithWorkingDirUsage("Working directory")
}

func GetWorkingDir() string {
	if !libs.IsValidPathArg(workingDir) || workingDir == "." {
		workingDir = libs.GetCurrentDir(true)
	}
	return workingDir
}

func WithRunArgs(runHandler func(...string)) {
	FuncRun = func(cmd *cobra.Command, args []string) {
		withCobraArgs(args)
		runHandler(args...)
	}
	setFuncs()
}

func WithRunArgsStr(runHandler func(string)) {
	FuncRun = func(cmd *cobra.Command, args []string) {
		withCobraArgs(args)
		runHandler(GetArgsStr())
	}
	setFuncs()
}

func WithRun(runHandler func()) {
	WithRunArgs(func(args ...string) {
		runHandler()
	})
}

func GetArgs() []string {
	return cobraArgs
}

func GetArgsStr() string {
	return slice.ArrayToString(GetArgs())
}

func DisableFlagParsing() {
	CobraCmd.DisableFlagParsing = true
}

func EnableFlagParsing() {
	CobraCmd.DisableFlagParsing = false
}
