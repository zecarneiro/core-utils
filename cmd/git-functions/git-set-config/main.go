package main

import (
	"fmt"
	"strings"

	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

const (
	globalFlags = "--global"
	cmdGit      = "git"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-set-config",
		Short: "Set git configuration (local and global)",
	}
	cobralib.WithRun(process)
}

func runCmd(cmdStr string) {
	exe.ExecRealTime(models.Command{Cmd: cmdStr, Verbose: true})
}

func runGitConfig(args string) {
	runCmd(fmt.Sprintf("%s %s", cmdGit, args))
}

func hasSubmodules(repoPath string) bool {
	output, err := exe.Exec(models.Command{Cmd: cmdGit, Args: []string{"submodule", "status"}, Cwd: repoPath})
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(output)) > 0
}

func process() {
	localCmds := []string{
		"config %s core.autocrlf input",
		"config %s core.fileMode false",
		"config %s core.logAllRefUpdates true",
		"config %s core.ignorecase true",
		"config %s pull.rebase true",
		"config %s --unset safe.directory",
		"config %s --add safe.directory '*'",
		"config %s merge.ff false",
	}
	globalCmds := append(localCmds, []string{
		"config %s --unset-all credential.helper",
		"config %s --unset credential.credentialstore",
	}...)
	currentDir, err := file.GetCurrentDir()
	// Process configs
	if platform.IsLinux() {
		if system.IsDesktopEnv(enums.KdeDE) {
			globalCmds = append(globalCmds, "config %s credential.credentialStore kwallet")
		} else {
			globalCmds = append(globalCmds, "config %s credential.credentialStore secretservice")
		}
	}
	for _, cmd := range globalCmds {
		runGitConfig(fmt.Sprintf(cmd, globalFlags))
	}
	logic.ProcessError(err)
	if file.IsDir(file.JoinPath(currentDir, ".git")) {
		logger.Header("Set local configurations")
		isRepoHasSubmodules := hasSubmodules(currentDir)
		for _, cmd := range localCmds {
			cmd = fmt.Sprintf(cmd, "")
			runGitConfig(cmd)
			if isRepoHasSubmodules {
				runGitConfig(fmt.Sprintf("submodule foreach --recursive \"%s %s\"", cmdGit, cmd))
			}
		}
	}
	runCmd("git-credential-manager configure")
}

func main() {
	cobralib.Run()
}
