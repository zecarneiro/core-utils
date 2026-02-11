package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
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
	cmd := fmt.Sprintf("sudo flatpak override %s --filesystem=host", appId)
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmd, Verbose: true}))
}

func main() {
	cobralib.Run()
}
