package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"main/internal/libs/cobralib"
	"strings"

	"github.com/spf13/cobra"
)

var (
	filter     string
	includeAll bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "snap-connect-all <package>",
		Short: "Enable all connections for a given package.",
	}
	cobralib.CobraCmd.Flags().StringVarP(&filter, "filter", "f", "", "Package to search")
	cobralib.CobraCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "Include all informations")
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	cmdInfo := models.Command{
		Cmd:      fmt.Sprintf(`snap connections --all | awk '$2 ~ "^%s:" && $3 == "-" {print $2}'`, app),
		Verbose:  false,
		UseShell: true,
	}
	cmdRes, err := exe.Exec(cmdInfo)
	logic.ProcessError(err)
	for connectionLine := range strings.SplitSeq(cmdRes, common.Eol()) {
		if !str.IsEmpty(connectionLine) {
			err := exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf("sudo snap connect '%s'", connectionLine), Verbose: true, UseShell: true})
			logger.Error(err)
		}
	}
}

func main() {
	cobralib.Run()
}
