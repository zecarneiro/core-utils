package main

import (
	"fmt"

	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	username string
	token    string
	repo     string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "github-change-url",
		Short: "Change GitHub remote URL with new username, token and repository",
	}
	cobralib.CobraCmd.Flags().StringVarP(&username, "username", "u", "", "Github username (required)")
	cobralib.CobraCmd.Flags().StringVarP(&token, "token", "t", "", "Github token (required)")
	cobralib.CobraCmd.Flags().StringVarP(&repo, "repository", "r", "", "Github repository (required)")
	cobralib.CobraCmd.MarkFlagRequired("username")
	cobralib.CobraCmd.MarkFlagRequired("token")
	cobralib.CobraCmd.MarkFlagRequired("repository")
	cobralib.WithRun(process)
}

func process() {
	fullURL := fmt.Sprintf("https://%s:%s@github.com/%s", username, token, repo)
	logger.Info(fmt.Sprintf("Will be set new github URL: %s", fullURL))
	cmd := models.Command{
		Cmd:     "git",
		Args:    []string{"remote", "set-url", "origin", fullURL},
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
