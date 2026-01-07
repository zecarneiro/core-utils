package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var filter string

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "snap-list",
		Short: "Get list of all packages or by filter",
	}
	cobralib.CobraCmd.Flags().StringVarP(&filter, "filter", "f", "", "Package to search")
	cobralib.WithRun(process)
}

func process() {
	cmdInfo := models.Command{Cmd: "snap list | awk 'NR >=2{print $1}'", UseShell: true}
	if len(filter) == 0 {
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	} else {
		cmdRes, err := exe.Exec(cmdInfo)
		logic.ProcessError(err)
		packagesList := strings.Split(cmdRes, common.Eol())
		for _, packageLine := range packagesList {
			if str.StringContains(packageLine, filter, true) {
				fmt.Println(packageLine)
			}
		}
	}
}

func main() {
	cobralib.Run()
}
