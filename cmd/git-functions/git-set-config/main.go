package main

import (
	"fmt"

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

const globalFlags = "--global"

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
	runCmd(fmt.Sprintf("git %s", args))
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
		for _, cmd := range localCmds {
			runGitConfig(fmt.Sprintf(cmd, ""))
		}
	}
	runCmd("git-credential-manager configure")
}

func main() {
	cobralib.Run()
}
