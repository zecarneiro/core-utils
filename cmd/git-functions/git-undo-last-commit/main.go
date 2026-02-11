package main

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-undo-last-commit",
		Short: "Undo the last commit (soft reset to HEAD~1)",
	}
	cobralib.WithRun(process)
}

func process() {
	cmd := models.Command{
		Cmd:     "git",
		Args:    []string{"reset", "--soft", "HEAD~1"},
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
