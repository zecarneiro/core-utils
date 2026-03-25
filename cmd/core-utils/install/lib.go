package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/system"
	"main/internal/dir"
	"main/internal/libs"
)

func getCmdInfo() models.Command {
	cmdInfo := models.Command{
		Verbose:  true,
		UseShell: true,
		IsThrow:  true,
	}
	if platform.IsLinux() && shell.IsBashInstalled() {
		cmdInfo.ShellToUse = enums.Bash
	} else if platform.IsWindows() && shell.IsPowershellInstalled() && shell.IsPromptCMDInstalled() {
		cmdInfo.ShellToUse = enums.PowerShell
	} else {
		logic.ProcessError(errors.New(common.NotImplementedYetMSG))
	}
	return cmdInfo
}

func resetCmdInfo(cmdInfo models.Command) models.Command {
	cmdInfo.Cmd = ""
	cmdInfo.Args = []string{}
	return cmdInfo
}

func askProcessPackage(packageName string) bool {
	return console.Confirm(fmt.Sprintf("Do you want to process all state for %s", packageName), true)
}

func addCoreUtilsDirsOnPathEnv() {
	envManager.Sync(envPathName)
	envPathValues := envManager.GetEnvValues(envPathName)
	dirList := []string{
		dir.CoreUtilsPrompt(),
		dir.CoreUtilsShellScripts(),
		dir.CoreUtilsSystemInstallBin(),
	}
	for _, dirName := range dirList {
		logger.Info(fmt.Sprintf("Will be Inserted '%s' on %s env", dirName, envPathName))
		envPathValues = append([]string{dirName}, envPathValues...)
	}
	envManager.UpdateEnv(envPathName, envPathValues)
	envManager.Sync(envPathName)
}

func addUserBinOnPathEnv() {
	envManager.Sync(envPathName)
	binDir := system.HomeUserBinDir()
	logic.ProcessError(file.CreateDirectory(binDir, true))
	envPathValues := envManager.GetEnvValues(envPathName)
	logger.Info(fmt.Sprintf("Insert '%s' on %s env", binDir, envPathName))
	envPathValues = append([]string{binDir}, envPathValues...)
	envManager.UpdateEnv(envPathName, envPathValues)
	envManager.Sync(envPathName)
}

func installShellScriptOnSystemProfile(fileScript string, data string) {
	canInsert := true
	if file.IsFile(fileScript) {
		exist, err := file.FileTextContains(fileScript, data, false)
		logic.ProcessError(err)
		if exist {
			canInsert = false
		}
	}
	if canInsert {
		fileConfig := models.FileWriterConfig{
			File:        fileScript,
			Data:        data,
			IsAppend:    true,
			WithUtf8BOM: false,
			IsCreateDir: true,
		}
		logic.ProcessError(file.WriteFile(fileConfig))
		logger.Ok(fmt.Sprintf("Update %s shell profile script with coreutils script", fileScript))
	}
}

func cleanEnvPath() {
	envManager.Sync(envPathName)
	envPathValues := envManager.GetEnvValues(envPathName)
	envPathValues = envManager.RemoveDuplicated(envPathValues)
	envManager.UpdateEnv(envPathName, envPathValues)
}

func createDirs() {
	logger.Header("Create System Dirs")
	cmds := []string{
		"CONFIG_DIR",
		"OTHER_APPS_DIR",
		"TEMP_DIR",
		"USER_BIN_DIR",
		"USER_STARTUP_DIR",
		"USER_TEMP_DIR",
	}
	for _, cmd := range cmds {
		libs.RunCoreUtilsCmd(cmd, false)
	}
}
