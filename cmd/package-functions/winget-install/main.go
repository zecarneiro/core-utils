package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	sourceArg        string
	isInteractiveArg bool
	validSource      = []string{"winget", "msstore"}
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "winget-install <app_id>",
		Short: "Install Winget app by given app id",
	}
	cobralib.CobraCmd.Flags().StringVarP(&sourceArg, "source", "s", validSource[0], fmt.Sprintf("Source(repository) to use. Valid value: %s. (Defalult: %s)", slice.ArrayToStringBySep(validSource, ","), validSource[0]))
	cobralib.CobraCmd.Flags().BoolVarP(&isInteractiveArg, "interactive", "i", false, "Enable Interactive instalation")
	cobralib.WithRunArgsStr(process)
}

func process(appId string) {
	if str.IsEmpty(appId) {
		logic.ProcessError(errors.New("invalid given app id"))
	}
	interactiveFlag := logic.Ternary(isInteractiveArg, "--interactive", "")
	cmd := fmt.Sprintf("winget install %s --accept-source-agreements --accept-package-agreements --source %s %s", appId, sourceArg, interactiveFlag)
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmd, Verbose: true}))
}

func main() {
	cobralib.Run()
}
