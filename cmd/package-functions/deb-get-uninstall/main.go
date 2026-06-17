package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "deb-get-uninstall <app>",
		Short: "Uninstall deb-get app",
		Args:  cobra.ExactArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	cmdStr := fmt.Sprintf("sudo deb-get purge %s", app)
	if console.CmdExists("deb-get") {
		logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: true}))
	} else {
		logger.Warn("DEB-GET not found on your system")
	}
}

func main() {
	cobralib.Run()
}
