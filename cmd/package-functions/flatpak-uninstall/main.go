package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var withSudo bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "flatpak-uninstall <app>",
		Short: "Uninstall flatpak app",
		Args:  cobra.ExactArgs(1),
	}
	cobralib.CobraCmd.Flags().BoolVarP(&withSudo, "with-sudo", "S", false, "Run with sudo")
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	cmdStr := fmt.Sprintf("flatpak uninstall --delete-data -y %s", app)
	if withSudo {
		cmdStr = fmt.Sprintf(`sudo %s`, cmdStr)
	}
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: true}))
}

func main() {
	cobralib.Run()
}
