package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	packageName    string
	othersSnapArgs []string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "snap-install",
		Short: "Install package and enable all connections for a given package.",
	}
	cobralib.CobraCmd.Flags().StringVarP(&packageName, "package", "p", "", "Package to install")
	cobralib.CobraCmd.Flags().StringArrayVarP(&othersSnapArgs, "others-snap-args", "a", make([]string, 0), "Others snap args.")
	cobralib.CobraCmd.MarkFlagRequired("package")
	cobralib.WithRun(process)
}

func process() {
	cmdStr := fmt.Sprintf(`sudo snap install %s`, packageName)
	if str.IsEmpty(packageName) {
		logic.ProcessError(fmt.Errorf("Invalid given package."))
	}
	if len(othersSnapArgs) > 0 {
		cmdStr = fmt.Sprintf(`%s %s`, cmdStr, slice.ArrayToString(othersSnapArgs))
	}
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: true, UseShell: true}))
	libs.RunCoreUtilsCmd("snap-connect-all", false, packageName)
}

func main() {
	cobralib.Run()
}
