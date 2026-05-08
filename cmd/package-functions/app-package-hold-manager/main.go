package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"slices"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	packageArg    string
	isUnholdArg   bool
	validPackages []string
)

func init() {
	loadData()
	setupCommand()
}

func loadData() {
	if platform.IsWindows() {
		validPackages = []string{"winget", "scoop"}
	} else if platform.IsLinux() {
		validPackages = []string{"apt", "flatpak", "snap"}
	} else {
		validPackages = []string{}
	}
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "app-package-hold-manager <app-id>",
		Short: "Avoid app to upgrade from package manager",
	}
	cobralib.CobraCmd.Flags().StringVarP(&packageArg, "package-manager", "p", "", fmt.Sprintf("Package manager which app is from. Accept only %s", slice.ArrayToStringBySep(validPackages, ",")))
	cobralib.CobraCmd.Flags().BoolVarP(&isUnholdArg, "unhold", "u", false, "UnHold APP")
	cobralib.WithRunArgsStr(process)
}

func validate(appArg string) {
	if str.IsEmpty(appArg) {
		logic.ProcessError(fmt.Errorf("APP can not be empty."))
	}
	if !slices.Contains(validPackages, packageArg) {
		logic.ProcessError(fmt.Errorf("Invalid given package manager. Accept only: %s", slice.ArrayToStringBySep(validPackages, ",")))
	}
}

func process(appArg string) {
	cmdInfo := models.Command{Cmd: "", Verbose: true, UseShell: true}
	validate(appArg)
	switch packageArg {
	case "winget":
		operationFlag := logic.Ternary(isUnholdArg, "remove", "add")
		cmdInfo.Cmd = fmt.Sprintf(`winget pin %s --id %s`, operationFlag, appArg)
	case "scoop":
		operationFlag := logic.Ternary(isUnholdArg, "unhold", "hold")
		cmdInfo.Cmd = fmt.Sprintf(`scoop %s %s`, operationFlag, appArg)
	case "apt":
		operationFlag := logic.Ternary(isUnholdArg, "unhold", "hold")
		cmdInfo.Cmd = fmt.Sprintf(`sudo apt-mark %s %s`, operationFlag, appArg)
	case "flatpak":
		operationFlag := logic.Ternary(isUnholdArg, "--remove", "")
		cmdInfo.Cmd = fmt.Sprintf(`flatpak mask %s %s`, operationFlag, appArg)
	case "snap":
		operationFlag := logic.Ternary(isUnholdArg, "--unhold", "--hold=forever")
		cmdInfo.Cmd = fmt.Sprintf(`sudo snap refresh %s %s`, operationFlag, appArg)
	}
	if !str.IsEmpty(cmdInfo.Cmd) {
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	}
}

func main() {
	cobralib.Run()
}
