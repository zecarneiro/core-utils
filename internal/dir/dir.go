package dir

import (
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/system"
)

/* ----------------------------- CORE UTILS DIR ----------------------------- */
func CoreUtilsSystemInstall() string {
	dir := file.JoinPath(system.HomeUserOptDir(), "coreutils")
	logic.ProcessError(file.CreateDirectory(dir, true))
	return dir
}

func CoreUtilsSystemInstallShellScripts() string {
	dir := file.JoinPath(CoreUtilsSystemInstall(), "scripts")
	logic.ProcessError(file.CreateDirectory(dir, true))
	return dir
}

func CoreUtilsSystemInstallBin() string {
	dir := file.JoinPath(CoreUtilsSystemInstall(), "bin")
	logic.ProcessError(file.CreateDirectory(dir, true))
	return dir
}

func CoreUtilsLocal() string {
	directory := file.JoinPath(system.HomeUserLocalDir(), "share", "coreutils")
	if !file.IsDir(directory) {
		logic.ProcessError(file.CreateDirectory(directory, true))
	}
	return directory
}

func CoreUtilsShellScripts() string {
	directory := file.JoinPath(CoreUtilsLocal(), "scripts")
	if !file.IsDir(directory) {
		logic.ProcessError(file.CreateDirectory(directory, true))
	}
	return directory
}

func CoreUtilsPrompt() string {
	directory := file.JoinPath(CoreUtilsLocal(), "prompt")
	if !file.IsDir(directory) {
		logic.ProcessError(file.CreateDirectory(directory, true))
	}
	return directory
}

func CoreUtilsUserConfig() string {
	directory := file.JoinPath(system.HomeUserConfigDir(), "coreutils")
	logic.ProcessError(file.CreateDirectory(directory, true))
	return directory
}
