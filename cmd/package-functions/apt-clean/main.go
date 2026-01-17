package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/entity"
	"golangutils/pkg/generic"
	"golangutils/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command

	aptCmdList = []string{
		"sudo apt clean -y",
		"sudo apt autoremove -y",
		"sudo apt autopurge -y",
		"sudo apt autoclean -y",
	}
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "apt-clean",
		Short: "Cleanup APT",
	}
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		logger.Title("Cleanup APT")
		for _, aptCmd := range aptCmdList {
			generic.ProcessError(console.ExecRealTime(entity.Command{Cmd: aptCmd, Verbose: true, UseShell: true}))
		}
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
