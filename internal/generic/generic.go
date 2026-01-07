package generic

import (
	"golangutils/pkg/common/platform"
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/system"
	"main/internal/dir"
)

func GetAllShellProfilesFiles() map[console.ShellType]string {
	shells := map[console.ShellType]string{
		console.Bash: file.ResolvePath(system.HomeDir(), ".bashrc"),
		console.Zsh:  file.ResolvePath(system.HomeDir(), ".zshrc"),
		console.Fish: file.ResolvePath(dir.GetUserConfig(), "fish/config.fish"),
		console.Ksh:  file.ResolvePath(system.HomeDir(), ".kshrc"),
	}
	if platform.IsWindows() {
		shells[console.PowerShell] = file.ResolvePath(system.HomeDir(), "Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1")
	} else if platform.IsPlatform([]platform.PlatformType{platform.Linux, platform.Darwin}) {
		shells[console.PowerShell] = file.ResolvePath(dir.GetUserConfig(), "powershell/Microsoft.PowerShell_profile.ps1")
	}
	return shells
}
