package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "flatpak-set-permission <app_id>",
		Short: "Set Flatpak permission for given app id",
	}
	cobralib.WithRunArgsStr(process)
}

func process(appId string) {
	if str.IsEmpty(appId) {
		logic.ProcessError(errors.New("invalid given app id"))
	}
	cmds := []string{
		"--filesystem=host",
		"--device=all",
		"--share=network",
		"--share=ipc",
		"--socket=pulseaudio",
	}
	for _, cmdSuffix := range cmds {
		permCmds := []string{
			fmt.Sprintf("flatpak --user override %s %s", appId, cmdSuffix),
			fmt.Sprintf("sudo flatpak --system override %s %s", appId, cmdSuffix),
		}
		for _, cmd := range permCmds {
			logger.Error(exe.ExecRealTime(models.Command{Cmd: cmd, Verbose: true}))
		}
	}
}

func main() {
	cobralib.Run()
}
