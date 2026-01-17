package dir

import (
	"golangutils/pkg/common"
	"golangutils/pkg/common/platform"
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"golangutils/pkg/system"
)

func GetUserConfig() string {
	config_dir := file.ResolvePath(system.HomeDir() + "/.config")
	generic.ProcessError(file.CreateDirectory(config_dir, true))
	return config_dir
}

func GetUserLocal() string {
	config_dir := file.ResolvePath(system.HomeDir() + "/.local")
	generic.ProcessError(file.CreateDirectory(config_dir, true))
	return config_dir
}

func GetUserOpt() string {
	opt_dir := file.ResolvePath(GetUserLocal() + "/opt")
	generic.ProcessError(file.CreateDirectory(opt_dir, true))
	return opt_dir
}

func GetUserBin() string {
	bin_dir := file.ResolvePath(GetUserLocal() + "/bin")
	generic.ProcessError(file.CreateDirectory(bin_dir, true))
	return bin_dir
}

func GetUserStartup() string {
	startup_dir := common.Unknown
	if platform.IsWindows() {
		startup_dir = file.ResolvePath(system.HomeDir() + "\\Start Menu\\Programs\\Startup")
	} else if platform.IsLinux() {
		startup_dir = file.ResolvePath(GetUserConfig() + "/autostart")
	}
	generic.ProcessError(file.CreateDirectory(startup_dir, true))
	return startup_dir
}

func GetUserTemp() string {
	temp_dir := common.Unknown
	if platform.IsWindows() {
		temp_dir = file.ResolvePath(system.TempDir())
	} else if platform.IsLinux() {
		temp_dir = file.ResolvePath(GetUserLocal() + "/tmp")
	}
	generic.ProcessError(file.CreateDirectory(temp_dir, true))
	return temp_dir
}

func GetTemp() string {
	temp_dir := system.TempDir()
	return temp_dir
}

func GetCoreUtilsLocal() string {
	dir := file.ResolvePath(GetUserLocal(), "coreutils")
	generic.ProcessError(file.CreateDirectory(dir, true))
	return dir
}

func GetCoreUtilsShellFunctions() string {
	currentShell := console.GetCurrentShell()
	if !currentShell.IsValid() {
		return ""
	}
	shellDirName := currentShell.String()
	if console.IsShell([]console.ShellType{console.Cmd, console.PowerShell}) {
		shellDirName = common.Ternary(currentShell.Equals(console.Cmd), console.PowerShell.String(), currentShell.String())
	}
	directory := file.ResolvePath(GetCoreUtilsLocal(), "functions", shellDirName+"-shell")
	if !file.IsDir(directory) {
		generic.ProcessError(file.CreateDirectory(directory, true))
	}
	if platform.IsWindows() && console.IsShell([]console.ShellType{console.Cmd, console.PowerShell}) {
		cmdDir := file.ResolvePath(directory, console.Cmd.String())
		generic.ProcessError(file.CreateDirectory(cmdDir, true))
	}
	return directory
}

func GetCoreUtilsShellAlias() string {
	if !console.GetCurrentShell().IsValid() {
		return ""
	}
	directory := file.ResolvePath(GetCoreUtilsLocal(), "alias")
	if !file.IsDir(directory) {
		generic.ProcessError(file.CreateDirectory(directory, true))
	}
	if platform.IsWindows() {
		cmdDir := file.ResolvePath(directory, console.Cmd.String())
		generic.ProcessError(file.CreateDirectory(cmdDir, true))
	}
	return directory
}
