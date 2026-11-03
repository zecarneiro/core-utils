package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"main/internal/libs/cobralib"
	"os"

	"github.com/spf13/cobra"
)

var (
	envPathArg string
)

func init() {
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "sudo-set-user-env-path",
		Short: "Concatenate env PATH of User with ROOT to use on sudo command.",
	}

	// Add
	addCmd := &cobra.Command{
		Use:   "add <env-path>",
		Short: "Concatenate USER env Path with ROOT",
		Run: func(cmd *cobra.Command, args []string) {
			envPathArg = slice.ArrayToString(args)
			process(true)
		},
	}

	// Delete
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete concatenation",
		Run: func(cmd *cobra.Command, args []string) {
			process(false)
		},
	}
	// Add all commands flags
	cobralib.CobraCmd.AddCommand(addCmd, deleteCmd)
}

func process(isAdd bool) {
	if !system.IsAdmin() {
		logic.ProcessError(errors.New(system.NeedAdminAccessMsg))
	}
	securePathFile := "/etc/sudoers.d/custom-secure-path-cu"
	writer := models.FileWriterConfig{IsAppend: false, IsCreateDir: true}
	if isAdd {
		if str.IsEmpty(envPathArg) {
			logic.ProcessError(fmt.Errorf("Invalid given Env PATH!"))
		}
		writer.File = securePathFile
		writer.Data = fmt.Sprintf("Defaults secure_path=\"%s:%s\"", envPathArg, os.Getenv("PATH"))
		logic.ProcessError(file.WriteFile(writer))
		logger.Ok(fmt.Sprintf("Create file with name: %s", securePathFile))
	} else {
		logic.ProcessError(file.DeleteFile(securePathFile))
		logger.Ok(fmt.Sprintf("Deleted file with name: %s", securePathFile))
	}
}

func main() {
	cobralib.Run()
}
