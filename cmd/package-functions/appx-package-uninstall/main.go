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
		Use:   "appx-package-uninstall <app>",
		Short: "Uninstall AppxPackage app",
		Args:  cobra.ExactArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	cmdStr := fmt.Sprintf(`sudo powershell.exe -Command "Get-AppxPackage -AllUsers "%s" | Remove-AppxPackage -AllUsers"`, app)
	err := exe.ExecRealTime(models.Command{Cmd: cmdStr})
	logic.ProcessError(err)
}

func main() {
	cobralib.Run()
}
