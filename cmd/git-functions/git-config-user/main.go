package main

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	username string
	email    string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-config-user",
		Short: "Set git user.name and user.email configuration",
	}
	cobralib.CobraCmd.Flags().StringVarP(&username, "username", "u", "", "Github username")
	cobralib.CobraCmd.Flags().StringVarP(&email, "email", "e", "", "Github email")
	cobralib.WithRun(process)
}

func process() {
	alreadyRun := false
	globalFlag := ""
	if file.FileExist(".git") {
		globalFlag = "--global"
	}
	baseArgs := []string{"config", globalFlag}
	if username != "" {
		cmd := models.Command{
			Cmd:     "git",
			Args:    append(baseArgs, "user.name", username),
			Verbose: true,
		}
		logic.ProcessError(exe.ExecRealTime(cmd))
		alreadyRun = true
	}

	if email != "" {
		cmd := models.Command{
			Cmd:     "git",
			Args:    append(baseArgs, "user.email", email),
			Verbose: true,
		}
		logic.ProcessError(exe.ExecRealTime(cmd))
		alreadyRun = true
	}
	if !alreadyRun {
		cobralib.CobraCmd.Help()
	}
}

func main() {
	cobralib.Run()
}
