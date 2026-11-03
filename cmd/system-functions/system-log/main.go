package main

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	command string
	args    []string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "system-log",
		Short: "Show system log in realtime.",
	}
	cobralib.WithRun(process)
}

func process() {
	cmdStr := "sudo tail -f /var/log/syslog"
	if platform.IsWindows() {
		// Not working as expected
		cmdStr = "while ($true) { Get-WinEvent -LogName System -MaxEvents 1; Start-Sleep -Seconds 2 }"
	}
	logger.Error(exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: true, UseShell: true}))
}

func main() {
	cobralib.Run()
}
