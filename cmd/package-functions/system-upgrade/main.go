package main

import (
	"fmt"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"
	"main/internal/libs/golangutilslib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "system-upgrade",
		Short: "Upgrade System packages",
	}
	cobralib.WithRun(process)
}

func process() {
	cmd := "topgrade --cleanup --allow-root --skip-notify --yes --disable helm uv deb_get"
	if !golangutilslib.FuncIsWindows() {
		cmd = fmt.Sprintf("%s powershell", cmd)
	}
	logic.ProcessError(golangutilslib.FuncExecRealTime(models.Command{Cmd: cmd, UseShell: true, Verbose: true}))
	for _, cmdDebGet := range []string{"sudo deb-get update", "sudo deb-get upgrade"} {
		logic.ProcessError(golangutilslib.FuncExecRealTime(models.Command{Cmd: cmdDebGet, UseShell: true, Verbose: true}))
	}
	// TODO: script_updater_processor(["run", "--all"])
}

func main() {
	cobralib.Run()
}
