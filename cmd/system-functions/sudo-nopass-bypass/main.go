package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	scriptFileArg string
	systemUserArg string
)

func init() {
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "sudo-nopass-bypass",
		Short: "Set file to run as admin without asking for sudo password.",
	}
	cobralib.CobraCmd.Flags().StringVarP(&scriptFileArg, "script", "s", "", "Script to sudo bypass")
	cobralib.CobraCmd.Flags().StringVarP(&systemUserArg, "user", "u", "", "Your logged user")
	cobralib.CobraCmd.MarkFlagRequired("script")
	cobralib.CobraCmd.MarkFlagRequired("user")
	cobralib.WithRun(process)
}

func process() {
	if !system.IsAdmin() {
		logic.ProcessError(errors.New(system.NeedAdminAccessMsg))
	}
	if str.IsEmpty(systemUserArg) || !file.IsFile(scriptFileArg) {
		logic.ProcessError(fmt.Errorf("Invalid given command!"))
	}
	sudoersFile := fmt.Sprintf("/etc/sudoers.d/99-custom-%s-%s-nopasswd", file.FileName(scriptFileArg), systemUserArg)
	writer := models.FileWriterConfig{
		File:        sudoersFile,
		Data:        fmt.Sprintf("%s ALL=(ALL) NOPASSWD: %s", systemUserArg, scriptFileArg),
		IsAppend:    false,
		IsCreateDir: true,
	}
	logic.ProcessError(file.WriteFile(writer))
	logger.Ok(fmt.Sprintf("Create file with name: %s", sudoersFile))
}

func main() {
	cobralib.Run()
}
