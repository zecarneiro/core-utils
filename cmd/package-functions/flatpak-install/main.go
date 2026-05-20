package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"slices"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	withSudoArg         bool
	appIdFromFileArg    string
	installationTypeArg string
	installationType    []string
)

const (
	userTypeInstalation   = "user"
	systemTypeInstalation = "system"
)

func init() {
	installationType = []string{userTypeInstalation, systemTypeInstalation}
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "flatpak-install <app_id_or_file>",
		Short: "Install Flatpak app and set Flatpak permission for given app id",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&withSudoArg, "with-sudo", "S", false, "Run with sudo")
	cobralib.CobraCmd.Flags().StringVarP(&appIdFromFileArg, "app-id-from-file", "i", "", "App ID, if you pass installer file as arg. Use only for set some permission. (Optional)")
	cobralib.CobraCmd.Flags().StringVarP(&installationTypeArg, "installation-type", "t", systemTypeInstalation, fmt.Sprintf(`Installation type. Only accept: %s`, slice.ArrayToStringBySep(installationType, ",")))
	cobralib.WithRunArgsStr(process)
}

func validate(appIdOrFile string) {
	if str.IsEmpty(appIdOrFile) {
		logic.ProcessError(errors.New("invalid given app id or installer file"))
	}
	if !slices.Contains(installationType, installationTypeArg) {
		logic.ProcessError(errors.New("invalid given type of installation"))
	}
}

func getCmdByAppArg(appArg string) (string, string) {
	if file.IsFile(appArg) {
		return fmt.Sprintf(`flatpak install --%s %s -y`, installationTypeArg, appArg), appIdFromFileArg
	}
	return fmt.Sprintf(`flatpak install --%s flathub %s -y`, installationTypeArg, appArg), appArg
}

func process(appIdOrFile string) {
	validate(appIdOrFile)
	cmd, appId := getCmdByAppArg(appIdOrFile)
	if withSudoArg && installationTypeArg != userTypeInstalation {
		cmd = fmt.Sprintf(`sudo %s`, cmd)
	}
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmd, Verbose: true}))
	if !str.IsEmpty(appId) {
		libs.RunCoreUtilsCmd("flatpak-set-permission", false, appId)
	}
}

func main() {
	cobralib.Run()
}
