package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/enums"
	"golangutils/pkg/env"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/shell"
	"golangutils/pkg/slice"
	"slices"

	"main/internal/dir"
	"main/internal/libs"
)

func showDeleteDirMsg(directory string) {
	logger.Info(fmt.Sprintf("Deleting dir: %s", directory))
}

func deleteCoreUtilsDirsOnPathEnv() {
	canUpdate := false
	envManager.Sync(envPathName)
	envPathValues := envManager.GetEnvValues(envPathName)
	dirList := []string{
		dir.CoreUtilsPrompt(),
		dir.CoreUtilsShellScripts(),
		dir.CoreUtilsSystemInstallBin(),
	}
	for _, dirName := range dirList {
		if slices.Contains(envPathValues, dirName) {
			logger.Info(fmt.Sprintf("Will be deleted '%s' on %s env", dirName, envPathName))
			envPathValues = slice.FilterArray(envPathValues, func(val string) bool {
				return val != dirName
			})
			canUpdate = true
		}
	}
	if canUpdate {
		envManager.UpdateEnv(envPathName, envPathValues)
		envManager.Sync(envPathName)
	}
}

func uninstallShellScriptOnSystemProfile(fileScript string, data string) {
	if file.IsFile(fileScript) {
		file.DeleteFileLines(fileScript, data, false)
		exist, err := file.FileTextContains(fileScript, data, false)
		logic.ProcessError(err)
		if !exist {
			logger.Ok(fmt.Sprintf("Delete from \"%s\" all data with coreutils script", fileScript))
		} else {
			logic.ProcessError(fmt.Errorf("Can't delete all data in system profile \"%s\" script with coreutils", fileScript))
		}
	}
}

func deleteAllNecessaryFiles() {
	coreUtilsLocalDir := dir.CoreUtilsLocal()
	coreUtilsUserConfigDir := dir.CoreUtilsUserConfig()
	showDeleteDirMsg(coreUtilsLocalDir)
	logic.ProcessError(file.DeleteDirectory(coreUtilsLocalDir))
	showDeleteDirMsg(coreUtilsUserConfigDir)
	logic.ProcessError(file.DeleteDirectory(coreUtilsUserConfigDir))
	showDeleteDirMsg(systemInstallDir)
	logic.ProcessError(file.DeleteDirectory(systemInstallDir))
}

func loadAll() {
	console.EnableFeatures()
	envPathName = env.GetPathName()
	envManager = libs.NewEnvManager()
	systemInstallDir = dir.CoreUtilsSystemInstall()
}

func main() {
	loadAll()
	logger.Title(uninstallTitleMsg)
	// Delete profiles scripts from system shell profile script's
	uninstallShellScriptOnSystemProfile(shell.GetShellProfileFile(enums.Bash), systemInstallDir)
	uninstallShellScriptOnSystemProfile(shell.GetShellProfileFile(enums.PowerShell), systemInstallDir)
	deleteCoreUtilsDirsOnPathEnv()
	deleteAllNecessaryFiles()
	logger.Ok("Uninstall done. Please, restart your terminal!")
}
