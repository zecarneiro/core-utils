package args

import (
	"errors"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

type GithubArgs struct {
	BaseArgs   // Extend
	Owner      string
	Repository string
	IsLatest   bool
}

var (
	ownerFlags = cobralib.FlagsOptions[string]{
		Name:       "owner",
		Shorthand:  "o",
		Usage:      "Owner of the repository",
		IsRequired: true,
		Default:    "",
	}
	repositoryFlags = cobralib.FlagsOptions[string]{
		Name:       "repository",
		Shorthand:  "r",
		Usage:      "Repository",
		IsRequired: true,
		Default:    "",
	}
	isLatestFlags = cobralib.FlagsOptions[bool]{
		Name:       "latest",
		Shorthand:  "l",
		Usage:      "Enable to get latest version",
		IsRequired: false,
		Default:    false,
	}
)

func BuildGithubFlags(inputCmd *cobra.Command) {
	buildBaseFlags(inputCmd)
	// Owner
	inputCmd.Flags().StringP(ownerFlags.Name, ownerFlags.Shorthand, ownerFlags.Default, ownerFlags.Usage)
	logic.ProcessError(inputCmd.MarkFlagRequired(ownerFlags.Name))
	// Repository
	inputCmd.Flags().StringP(repositoryFlags.Name, repositoryFlags.Shorthand, repositoryFlags.Default, repositoryFlags.Usage)
	logic.ProcessError(inputCmd.MarkFlagRequired(repositoryFlags.Name))
	// IsLatest
	inputCmd.Flags().BoolP(isLatestFlags.Name, isLatestFlags.Shorthand, isLatestFlags.Default, isLatestFlags.Usage)
}

func (g *GithubArgs) ValidateGithub() {
	g.validateBase()
	if str.IsEmpty(g.Owner) {
		logic.ProcessError(errors.New("invalid given owner"))
	}
	if str.IsEmpty(g.Repository) {
		logic.ProcessError(errors.New("invalid given repository"))
	}
}

func (g *GithubArgs) FillGithubValues(cmd *cobra.Command) {
	var err error
	g.fillBaseValues(cmd)
	// Owner
	g.Owner, err = cmd.Flags().GetString(ownerFlags.Name)
	logic.ProcessError(err)
	// Repository
	g.Repository, err = cmd.Flags().GetString(repositoryFlags.Name)
	logic.ProcessError(err)
	// Is Latest
	g.IsLatest, err = cmd.Flags().GetBool(isLatestFlags.Name)
	logic.ProcessError(err)
}
