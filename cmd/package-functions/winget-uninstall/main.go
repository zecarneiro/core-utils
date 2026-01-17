package main

import (
	"fmt"
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
		Use:   "winget-uninstall [app]",
		Short: "Uninstall winget app",
		Args:  cobra.MinimumNArgs(1),
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		data := strings.Join(args, " ")
		cmdStr := fmt.Sprintf("winget uninstall --purge %s", data)
		err := console.ExecRealTime(entity.Command{Cmd: cmdStr, Verbose: true})
		generic.ProcessError(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
