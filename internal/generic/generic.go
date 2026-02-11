package generic

import (
	"golangutils/pkg/file"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/system"

	"main/internal/dir"
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
