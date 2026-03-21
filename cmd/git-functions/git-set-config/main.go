package main

import (
	"strings"

	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-set-config",
		Short: "Set git configuration (local and global)",
	}
	cobralib.WithRun(process)
}

func process() {
	localCmds := []string{
		"config core.autocrlf input",
		"config core.fileMode false",
		"config core.logAllRefUpdates true",
		"config core.ignorecase true",
		"config pull.rebase true",
		"config --unset safe.directory",
		"config --add safe.directory '*'",
		"config merge.ff false",
	}
	globalCmds := make([]string, len(localCmds))
	for i, cmd := range localCmds {
		globalCmds[i] = "--global " + cmd
	}
	globalCmds = append(globalCmds, []string{
		"config --global credential.credentialStore plaintext",
		"git-credential-manager configure",
	}...)
	for _, cmd := range globalCmds {
		runGitConfig(cmd)
	}
	if file.IsDir(".git") {
		logger.Info("Set local configurations")
		for _, cmd := range localCmds {
			runGitConfig(cmd)
		}
	}
}

func runGitConfig(args string) {
	parts := strings.Fields(args)
	cmd := models.Command{
		Cmd:     "git",
		Args:    parts,
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
