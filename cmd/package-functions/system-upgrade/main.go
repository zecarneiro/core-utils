package main

import (
	"fmt"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"

	"main/internal/libs"
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
	} else {
		cmd = fmt.Sprintf("%s system", cmd)
	}
	golangutilslib.FuncExecRealTime(models.Command{Cmd: cmd, UseShell: true})
	if platform.IsLinux() {
		for _, cmdDebGet := range []string{"sudo deb-get update", "sudo deb-get upgrade"} {
			golangutilslib.FuncExecRealTime(models.Command{Cmd: cmdDebGet, Verbose: true, UseShell: true})
		}
	}
	libs.RunCoreUtilsCmd("script-manager-cu", false, "run")
}

func main() {
	cobralib.Run()
}
