package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "flatpak-install <app_id>",
		Short: "Install Flatpak app and Set Flatpak permission for given app id",
	}
	cobralib.WithRunArgsStr(process)
}

func process(appId string) {
	if str.IsEmpty(appId) {
		logic.ProcessError(errors.New("invalid given app id"))
	}
	cmd := fmt.Sprintf("flatpak install flathub %s -y", appId)
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmd, Verbose: true}))
	libs.RunCoreUtilsCmd("flatpak-set-permission", false, appId)
}

func main() {
	cobralib.Run()
}
