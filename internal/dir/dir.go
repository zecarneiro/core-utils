package dir

import (
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/system"
)

func GetUserConfig() string {
	config_dir := file.ResolvePath(system.HomeDir() + "/.config")
	logic.ProcessError(file.CreateDirectory(config_dir, true))
	return config_dir
}

func GetUserLocal() string {
	config_dir := file.ResolvePath(system.HomeDir() + "/.local")
	logic.ProcessError(file.CreateDirectory(config_dir, true))
	return config_dir
}

func GetUserOpt() string {
	opt_dir := file.ResolvePath(GetUserLocal() + "/opt")
	logic.ProcessError(file.CreateDirectory(opt_dir, true))
	return opt_dir
}

func GetUserBin() string {
	bin_dir := file.ResolvePath(GetUserLocal() + "/bin")
	logic.ProcessError(file.CreateDirectory(bin_dir, true))
	return bin_dir
}

func GetUserStartup() string {
	startup_dir := common.Unknown
	if platform.IsWindows() {
		startup_dir = file.ResolvePath(system.HomeDir() + "\\Start Menu\\Programs\\Startup")
	} else if platform.IsLinux() {
		startup_dir = file.ResolvePath(GetUserConfig() + "/autostart")
	}
	logic.ProcessError(file.CreateDirectory(startup_dir, true))
	return startup_dir
}

func GetUserTemp() string {
	temp_dir := common.Unknown
	if platform.IsWindows() {
		temp_dir = file.ResolvePath(system.TempDir())
	} else if platform.IsLinux() {
		temp_dir = file.ResolvePath(GetUserLocal() + "/tmp")
	}
	logic.ProcessError(file.CreateDirectory(temp_dir, true))
	return temp_dir
}

func GetTemp() string {
	temp_dir := system.TempDir()
	return temp_dir
}

func GetCoreUtilsLocal() string {
	dir := file.ResolvePath(GetUserLocal(), "coreutils")
	logic.ProcessError(file.CreateDirectory(dir, true))
	return dir
}

func GetCoreUtilsShellFunctions() string {
	currentShell := shell.GetCurrentShell()
	if !currentShell.IsValid() {
		return ""
	}
	shellDirName := currentShell.String()
	if shell.IsShell([]shell.ShellType{shell.Cmd, shell.PowerShell}) {
		shellDirName = logic.Ternary(currentShell.Equals(shell.Cmd), shell.PowerShell.String(), currentShell.String())
	}
	directory := file.ResolvePath(GetCoreUtilsLocal(), "functions", shellDirName+"-shell")
	if !file.IsDir(directory) {
		logic.ProcessError(file.CreateDirectory(directory, true))
	}
	if platform.IsWindows() && shell.IsShell([]shell.ShellType{shell.Cmd, shell.PowerShell}) {
		cmdDir := file.ResolvePath(directory, shell.Cmd.String())
		logic.ProcessError(file.CreateDirectory(cmdDir, true))
	}
	return directory
}

func GetCoreUtilsShellAlias() string {
	if !shell.GetCurrentShell().IsValid() {
		return ""
	}
	directory := file.ResolvePath(GetCoreUtilsLocal(), "alias")
	if !file.IsDir(directory) {
		logic.ProcessError(file.CreateDirectory(directory, true))
	}
	if platform.IsWindows() {
		cmdDir := file.ResolvePath(directory, shell.Cmd.String())
		logic.ProcessError(file.CreateDirectory(cmdDir, true))
	}
	return directory
}
