package libs

import (
	"golangutils/pkg/common"
	"golangutils/pkg/env"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/system"

	"main/internal/dir"
	"main/internal/libs/golangutilslib"
)

func GetAllShellProfilesFiles() map[shell.ShellType]string {
	shells := map[shell.ShellType]string{
		shell.Bash: file.ResolvePath(system.HomeDir(), ".bashrc"),
		shell.Zsh:  file.ResolvePath(system.HomeDir(), ".zshrc"),
		shell.Fish: file.ResolvePath(dir.GetUserConfig(), "fish/config.fish"),
		shell.Ksh:  file.ResolvePath(system.HomeDir(), ".kshrc"),
	}
	if platform.IsWindows() {
		shells[shell.PowerShell] = file.ResolvePath(system.HomeDir(), "Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1")
	} else if platform.IsPlatform([]platform.PlatformType{platform.Linux, platform.Darwin}) {
		shells[shell.PowerShell] = file.ResolvePath(dir.GetUserConfig(), "powershell/Microsoft.PowerShell_profile.ps1")
	}
	return shells
}

func GetCurrentDir(exitOnError bool) string {
	currentDir, err := golangutilslib.FuncGetCurrentDir()
	if exitOnError {
		logic.ProcessError(err)
	}
	return logic.Ternary(common.IsNil(err), currentDir, "")
}

func IsValidPathArg(path string) bool {
	return logic.Ternary(path == "" || path == " " || path == "..", false, true)
}

func RunCoreUtilsCmd(cmdName string, args ...string) {
	coreUtilExecDir, err := file.GetExecutableDir()
	logic.ProcessError(err)
	env.InsertOnPath(coreUtilExecDir)
	exe.ExecRealTime(models.Command{Cmd: cmdName, Args: args})
}
