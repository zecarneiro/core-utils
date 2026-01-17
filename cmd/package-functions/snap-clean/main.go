package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/entity"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"golangutils/pkg/logger"
	"golangutils/pkg/system"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "snap-clean",
		Short: "Cleanup SNAP",
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		snapCacheDir := "/var/lib/snapd/cache"
		logger.Title("Cleanup SNAP")
		if file.IsDir(snapCacheDir) {
			snapCacheCmd := fmt.Sprintf("rm -rf %s/*", snapCacheDir)
			generic.ProcessError(console.ExecRealTime(entity.Command{Cmd: fmt.Sprintf("sudo sh -c '%s'", snapCacheCmd), Verbose: true}))
		}
		logger.Warn("Removes old revisions of snaps")
		if generic.Confirm("Please, CLOSE ALL SNAPS BEFORE RUNNING THIS. Continue?", false) {
			cmdGetListApps := "LANG=en_US.UTF-8 snap list --all | awk '/disabled/{print $1, $3}'"
			cmdRes, err := console.Exec(entity.Command{Cmd: cmdGetListApps, UseShell: true})
			generic.ProcessError(err)
			cmdResList := strings.Split(cmdRes, system.Eol())
			for _, cmdRespData := range cmdResList {
				data := strings.Split(cmdRespData, " ")
				if len(data) == 2 {
					snapRemoveCmd := fmt.Sprintf("sudo snap remove \"%s\" --revision=\"%s\"", data[0], data[1])
					err := console.ExecRealTime(entity.Command{Cmd: snapRemoveCmd, Verbose: true, UseShell: true})
					if err != nil {
						logger.Error(err)
					}
				}
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
