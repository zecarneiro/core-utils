package main

import (
	"errors"
	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/env"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"

	"main/internal/dir"
	"main/internal/libs"
)

func copyAllNecessaryFiles() {
	startScriptName := logic.Ternary(platform.IsWindows(), "start.bat", "start.sh")
	logger.Info("Copy all files and directories to app system directory")
	logic.ProcessError(file.CopyDir(rootDir, systemInstallDir))
	logic.ProcessError(file.DeleteDirectory(file.JoinPath(systemInstallDir, "installers")))
	logic.ProcessError(file.DeleteFile(file.JoinPath(systemInstallDir, startScriptName)))
}

func isCoreUtilsInstallFolderEmpty() bool {
	if !file.IsDir(systemInstallDir) {
		return true
	}
	filesInfo, err := file.ReadDirRecursive(systemInstallDir)
	logic.ProcessError(err)
	return len(filesInfo.Files) == 0
}

func validate() {
	if !system.IsValidUserHomeDir(true) {
		logic.ProcessError(errors.New("Invalid Username/User home directory"))
	}
	if !(platform.IsLinux() && shell.IsBashInstalled()) && !(platform.IsWindows() && shell.IsPowershellInstalled() && shell.IsPromptCMDInstalled()) {
		logic.ProcessError(errors.New(common.NotImplementedYetMSG))
	}
	if str.IsEmpty(rootDir) || !file.IsDir(rootDir) {
		logic.ProcessError(errors.New("Please, pass a valid root dir from start.ps1 or start.sh"))
	}
}

func loadAndValidate() {
	isUpdateOnly = false
	console.EnableFeatures()
	libs.EnableRunningFromCUInstallDir()
	rootDir = slice.ArrayToString(console.GetArgsList())
	envManager = libs.NewEnvManager()
	envPathName = env.GetPathName()
	validate()
	envManager.SetSystemConfig()
	systemInstallDir = dir.CoreUtilsSystemInstall() // necessary to next lines
	if !isCoreUtilsInstallFolderEmpty() {
		logger.Info("Detected that CoreUtils is already installed on this system.")
		if console.Confirm("Will be process update CU only. Continue?", false) {
			isUpdateOnly = true
		} else {
			logic.ProcessError(errors.New("Please, uninstall before to continue."))
		}
	}
}

func main() {
	loadAndValidate()
	logger.Title(installTitleMsg)
	copyAllNecessaryFiles()
	if !platform.IsWindows() {
		libs.RunCoreUtilsCmd("777", true, systemInstallDir)
	}
	envManager.Sync(envPathName)
	if platform.IsLinux() {
		linuxProcessor := NewDependencyLinux(isUpdateOnly)
		linuxProcessor.start()
	} else if platform.IsWindows() {
		windowsProcessor := NewDependencyWindows(isUpdateOnly)
		windowsProcessor.start()
	}
	logger.Ok("Done.")
	if isUpdateOnly {
		logger.Warn("Please, restart your terminal.")
	} else {
		logger.Error(system.Reboot())
	}
}
