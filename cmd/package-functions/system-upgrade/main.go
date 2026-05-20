package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/logger"
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
	cmd := "topgrade --cleanup --allow-root --skip-notify --yes --disable helm uv deb_get containers"
	logger.WithIgnoreExitStatusError(true)
	if !golangutilslib.FuncIsWindows() {
		cmd = fmt.Sprintf("%s powershell", cmd)
	} else {
		cmd = fmt.Sprintf("%s system", cmd)
		logger.Error(golangutilslib.FuncExecRealTime(models.Command{Cmd: "winget upgrade topgrade-rs.topgrade --force", UseShell: true, Verbose: true}))
	}
	logger.Error(golangutilslib.FuncExecRealTime(models.Command{Cmd: cmd, UseShell: true}))
	if platform.IsLinux() && console.CmdExists("deb-get") {
		for _, cmdDebGet := range []string{"sudo deb-get update", "sudo deb-get upgrade"} {
			logger.Error(golangutilslib.FuncExecRealTime(models.Command{Cmd: cmdDebGet, Verbose: true, UseShell: true}))
		}
	}
	logger.Title("Run all script(s) to install/update/remove package(s)")
	libs.RunCoreUtilsCmd("script-manager-cu", false, "run")
	logger.Title("Update All Extensions from all VSCode Profiles")
	libs.RunCoreUtilsCmd("vscode-extension-upgrade", false)
	logger.Title("Install/Upgrade all installer by URL")
	libs.RunCoreUtilsCmd("installer-by-url", false, "run")
}

func main() {
	cobralib.Run()
}
