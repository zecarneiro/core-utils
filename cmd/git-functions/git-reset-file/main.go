package main

import (
	"fmt"

	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

const masterBranch = "origin/master"

var (
	fileArg   string
	branchArg string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-reset-file",
		Short: "Reset a file to a specific branch",
	}
	cobralib.CobraCmd.Flags().StringVarP(&fileArg, "file", "f", "", "Filepath to reset (required)")
	cobralib.CobraCmd.Flags().StringVarP(&branchArg, "branch", "b", masterBranch, "Name of branch you want to file reset to (default: origin/master)")
	cobralib.CobraCmd.MarkFlagRequired("file")
	cobralib.WithRun(process)
}

func process() {
	if !file.IsFile(fileArg) {
		logic.ProcessError(fmt.Errorf("invalid given file: %s", fileArg))
	}
	cmd := models.Command{
		Cmd:     "git",
		Args:    []string{"checkout", branchArg, fileArg},
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
