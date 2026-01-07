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
		Use:   "wsl-uninstall [distro]",
		Short: "Uninstall wsl distro",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(distro string) {
	cmdStr := fmt.Sprintf("wsl --unregister %s", distro)
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: true}))
}

func main() {
	cobralib.Run()
}
