package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	resolveType  string
	validChoices = map[string]bool{"wsl2win": true, "win2wsl": true}
)

func init() { setupCommand() }

func setupCommand() {
	typeKey := "type"
	cobralib.CobraCmd = &cobra.Command{
		Use:   "resolve-wsl-path [path]",
		Short: "Resolve WSL path",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().StringVarP(&resolveType, typeKey, "t", "", "Type of path to resolve")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired(typeKey))
	cobralib.WithRunArgsStr(process)
}

func validateType() {
	keysArr := slice.MapToKeys(validChoices)
	if !validChoices[resolveType] {
		logic.ProcessError(fmt.Errorf("invalid type: %s. Valid choices are: %s", resolveType, strings.Join(keysArr, ", ")))
	}
}

func getCommand(path string) string {
	command := "wsl -- wslpath %s %s"
	switch resolveType {
	case "wsl2win":
		command = fmt.Sprintf(command, "-w", path)
	case "win2wsl":
		command = fmt.Sprintf(command, "-u", path)
	default:
		command = ""
	}
	return command
}

func process(path string) {
	validateType()
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: getCommand(path), Verbose: false, UseShell: true}))
}

func main() {
	cobralib.Run()
}
