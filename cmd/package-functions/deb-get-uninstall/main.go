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
		Use:   "deb-get-uninstall [app]",
		Short: "Uninstall deb-get app",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	cmdStr := fmt.Sprintf("sudo deb-get purge %s", app)
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: true}))
}

func main() {
	cobralib.Run()
}
