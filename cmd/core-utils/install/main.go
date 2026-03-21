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
	if !(platform.IsLinux() && shell.IsBashInstalled()) && !(platform.IsWindows() && shell.IsPowershellInstalled() && shell.IsPromptCMDInstalled()) {
		logic.ProcessError(errors.New(common.NotImplementedYetMSG))
	}
	if str.IsEmpty(rootDir) || !file.IsDir(rootDir) {
		logic.ProcessError(errors.New("Please, pass a valid root dir from start.ps1 or start.sh"))
	}
}

func loadAndValidate() {
	console.EnableFeatures()
	libs.EnableRunningFromCUInstallDir()
	rootDir = slice.ArrayToString(console.GetArgsList())
	envManager = libs.NewEnvManager()
	envPathName = env.GetPathName()
	validate()
	envManager.SetSystemConfig()
	systemInstallDir = dir.CoreUtilsSystemInstall() // necessary to next lines
	if !isCoreUtilsInstallFolderEmpty() {
		logic.ProcessError(errors.New("Detected that CoreUtils is already installed on this system. Please, uninstall before to continue."))
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
		linuxProcessor := NewDependencyLinux()
		linuxProcessor.start()
	} else if platform.IsWindows() {
		windowsProcessor := NewDependencyWindows()
		windowsProcessor.start()
	}
	logger.Ok("Install Done. Please, restart your terminal")
}
