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
		Use:   "git-cherrypick-master-continue",
		Short: "Continue a cherry-pick operation",
	}
	cobralib.WithRun(process)
}

func process() {
	cmd := models.Command{
		Cmd:     "git",
		Args:    []string{"cherry-pick", "--continue"},
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
