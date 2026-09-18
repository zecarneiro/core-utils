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
	}
	for _, cmdSuffix := range cmds {
		cmd := fmt.Sprintf("sudo flatpak override %s %s", appId, cmdSuffix)
		logger.Error(exe.ExecRealTime(models.Command{Cmd: cmd, Verbose: true}))
	}
}

func main() {
	cobralib.Run()
}
