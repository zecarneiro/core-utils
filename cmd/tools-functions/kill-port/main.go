package main

import (
	"fmt"

	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var portArg string

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "kill-port <port>",
		Short: "Kill a process running on a specific port",
		Args:  cobra.ExactArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(port string) {
	if !platform.IsWindows() && !platform.IsLinux() {
		logic.ProcessError(fmt.Errorf("unsupported platform"))
	}
	var cmd string
	if platform.IsLinux() {
		cmd = fmt.Sprintf("sudo kill -9 $(sudo lsof -t -i :%s)", port)
	} else {
		cmd = fmt.Sprintf("Get-Process -Id (Get-NetTCPConnection -LocalPort %s).OwningProcess | Stop-Process -Force", port)
	}
	execCmd := models.Command{
		Cmd:      cmd,
		Verbose:  true,
		UseShell: true,
	}
	logic.ProcessError(exe.ExecRealTime(execCmd))
}

func main() {
	cobralib.Run()
}
