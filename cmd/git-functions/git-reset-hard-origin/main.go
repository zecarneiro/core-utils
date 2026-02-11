package main

import (
	"fmt"
	"strings"

	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

const gitCmd = "git"

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-reset-hard-origin",
		Short: "Reset git hard to origin",
	}
	cobralib.WithRun(process)
}

func process() {
	cmd := models.Command{
		Cmd:     gitCmd,
		Args:    []string{"branch", "--show-current"},
		Verbose: true,
	}
	response, err := exe.Exec(cmd)
	logic.ProcessError(err)
	branchName := strings.TrimSpace(response)
	if str.IsEmpty(branchName) {
		logic.ProcessError(fmt.Errorf("could not get current branch name"))
	}
	resetCmd := models.Command{
		Cmd:     gitCmd,
		Args:    []string{"reset", "--hard", fmt.Sprintf("origin/%s", branchName)},
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(resetCmd))
}

func main() {
	cobralib.Run()
}
