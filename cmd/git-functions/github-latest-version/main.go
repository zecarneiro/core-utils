package main

import (
	"golangutils/pkg/git"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	owner      string
	repo       string
	latestArg  bool
	verboseArg bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "github-latest-version",
		Short: "Get latest version of a GitHub repository",
	}
	cobralib.CobraCmd.Flags().StringVarP(&owner, "owner", "o", "", "Github owner (required)")
	cobralib.CobraCmd.Flags().StringVarP(&repo, "repository", "r", "", "Github repository (required)")
	cobralib.CobraCmd.Flags().BoolVarP(&latestArg, "latest", "l", false, "Enable to get latest version")
	cobralib.CobraCmd.Flags().BoolVarP(&verboseArg, "verbose", "v", false, "Enable log on the screen")
	cobralib.CobraCmd.MarkFlagRequired("owner")
	cobralib.CobraCmd.MarkFlagRequired("repository")
	cobralib.WithRun(process)
}

func process() {
	noReleaseFoundMsg := "No releases found for this repository."
	alreadyRun := false
	release, err := git.GithubGetLatestVersionRepo(owner, repo, latestArg, verboseArg)
	logic.ProcessError(err)
	if !str.IsEmpty(release.Version) {
		logger.Log(release.Version)
		alreadyRun = true
	}
	if !alreadyRun {
		logger.Warn(noReleaseFoundMsg)
	}
}

func main() {
	cobralib.Run()
}
