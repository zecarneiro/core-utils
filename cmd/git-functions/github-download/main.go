package main

import (
	"errors"
	"golangutils/pkg/git"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	owner    string
	repo     string
	version  string
	filename string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "github-download",
		Short: "Download a file version from a GitHub repository",
	}
	cobralib.CobraCmd.Flags().StringVarP(&owner, "owner", "o", "", "Github owner (required)")
	cobralib.CobraCmd.Flags().StringVarP(&repo, "repository", "r", "", "Github repository (required)")
	cobralib.CobraCmd.Flags().StringVarP(&version, "version", "v", "", "Version of file (required)")
	cobralib.CobraCmd.Flags().StringVarP(&filename, "filename", "f", "", "Filename to download (required)")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("owner"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("repository"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("version"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("filename"))
	cobralib.WithRun(process)
}

func process() {
	if str.IsEmpty(owner) || str.IsEmpty(repo) || str.IsEmpty(version) || str.IsEmpty(filename) {
		logic.ProcessError(cobralib.CobraCmd.Help())
		logic.ProcessError(errors.New("Please, provide the required arguments"))
	}
	logic.ProcessError(git.DownloadFromGithubRepo(owner, repo, version, filename))
}

func main() {
	cobralib.Run()
}
