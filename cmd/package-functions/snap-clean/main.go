package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "snap-clean",
		Short: "Cleanup SNAP",
	}
	cobralib.WithRun(process)
}

func process() {
	logger.Title("Cleanup SNAP")
	logger.Warn("Removes old revisions of snaps")
	if console.Confirm("Please, CLOSE ALL SNAPS BEFORE RUNNING THIS. Continue?", false) {
		cmdGetListApps := "LANG=en_US.UTF-8 snap list --all | awk '/disabled/{print $1, $3}'"
		cmdRes, err := exe.Exec(models.Command{Cmd: cmdGetListApps, UseShell: true})
		logic.ProcessError(err)
		for cmdRespData := range strings.SplitSeq(cmdRes, common.Eol()) {
			data := strings.Split(cmdRespData, " ")
			if len(data) == 2 {
				snapRemoveCmd := fmt.Sprintf("sudo snap remove \"%s\" --revision=\"%s\"", data[0], data[1])
				err := exe.ExecRealTime(models.Command{Cmd: snapRemoveCmd, Verbose: true})
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}
}

func main() {
	cobralib.Run()
}
