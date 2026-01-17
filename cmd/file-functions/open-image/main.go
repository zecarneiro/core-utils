package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/common/platform"
	"golangutils/pkg/console"
	"golangutils/pkg/entity"
	"golangutils/pkg/generic"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "open-image [filepath]",
		Short: "Open image",
		Args:  cobra.MinimumNArgs(1),
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		var cmdInfo entity.Command
		data := strings.Join(args, " ")
		data = common.Ternary(data == "." || data == "..", "", data)
		if platform.IsWindows() {
			cmdInfo = entity.Command{Cmd: "Start-Process", Args: []string{data}, UseShell: true}
		} else if platform.IsLinux() {
			cmdInfo = entity.Command{Cmd: "xdg-open", Args: []string{data}, UseShell: false}
		} else if platform.IsDarwin() {
			cmdInfo = entity.Command{Cmd: "open", Args: []string{data}, UseShell: true}
		}
		cmdInfo.Verbose = false
		generic.ProcessError(console.ExecRealTime(cmdInfo))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
