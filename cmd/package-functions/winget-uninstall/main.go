package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "winget-uninstall [app]",
		Short: "Uninstall winget app",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	cmdStr := fmt.Sprintf("winget uninstall --purge %s", app)
	err := exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: true})
	logic.ProcessError(err)
}

func main() {
	cobralib.Run()
}
