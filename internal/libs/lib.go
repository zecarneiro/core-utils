package libs

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"golangutils/pkg/env"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/str"
	"strings"

	"main/internal/dir"
	"main/internal/libs/golangutilslib"
)

var runningFromCoreUtilsInstallDir = false // This var enable to run cmd from core utils install bin folder if activated

func EnableRunningFromCUInstallDir() {
	runningFromCoreUtilsInstallDir = true
}

func GetCurrentDir(exitOnError bool) string {
	currentDir, err := golangutilslib.FuncGetCurrentDir()
	if exitOnError {
		logic.ProcessError(err)
	}
	return logic.Ternary(common.IsNil(err), currentDir, "")
}

func IsValidPathArg(path string) bool {
	return logic.Ternary(str.IsEmpty(path) || path == "..", false, true)
}

func RunCoreUtilsCmdWithShell(cmdName string, verbose bool, shellToUse enums.ShellType, args ...string) {
	coreUtilExecDir, err := exe.GetExecutableDir()
	logic.ProcessError(err)
	if runningFromCoreUtilsInstallDir {
		coreUtilExecDir = dir.CoreUtilsSystemInstallBin()
	}
	useShell := logic.Ternary(shellToUse.IsValid(), true, false)
	env.InsertOnPath(coreUtilExecDir)
	logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmdName, Args: args, Verbose: verbose, UseShell: useShell, ShellToUse: shellToUse}))
}

func RunCoreUtilsCmd(cmdName string, verbose bool, args ...string) {
	RunCoreUtilsCmdWithShell(cmdName, verbose, enums.UnknownShell, args...)
}

func RunCoreUtilsCmdWithOutputWithShell(cmdName string, verbose bool, shellToUse enums.ShellType, args ...string) string {
	coreUtilExecDir, err := exe.GetExecutableDir()
	logic.ProcessError(err)
	if runningFromCoreUtilsInstallDir {
		coreUtilExecDir = dir.CoreUtilsSystemInstallBin()
	}
	useShell := logic.Ternary(shellToUse.IsValid(), true, false)
	env.InsertOnPath(coreUtilExecDir)
	output, err := exe.Exec(models.Command{Cmd: cmdName, Args: args, Verbose: verbose, UseShell: useShell, ShellToUse: shellToUse})
	logic.ProcessError(err)
	return strings.TrimSuffix(output, common.Eol())
}

func RunCoreUtilsCmdWithOutput(cmdName string, verbose bool, args ...string) string {
	return RunCoreUtilsCmdWithOutputWithShell(cmdName, verbose, enums.UnknownShell, args...)
}

func CreateExecPwshFromPromptCMD(pwshScriptFilepath string) {
	if platform.IsWindows() {
		filePrompt := file.JoinPath(dir.CoreUtilsLocal(), "prompt", fmt.Sprintf("%s.cmd", file.FileName(pwshScriptFilepath)))
		baseCmd := fmt.Sprintf("%s -NoLogo -ExecutionPolicy Bypass -File \"%s\" %s", shell.GetPowershellCmd(), pwshScriptFilepath, shell.CmdAllArgsVarStr)
		fileConfig := models.FileWriterConfig{
			File:        filePrompt,
			Data:        fmt.Sprintf(`@echo off%s%s%s`, common.Eol(), baseCmd, common.Eol()),
			IsAppend:    false,
			IsCreateDir: true,
			WithUtf8BOM: false,
		}
		logic.ProcessError(file.WriteFile(fileConfig))
	}
}

func GetScriptLibPathByName(name string) string {
	return file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "libs", name)
}

func GetScriptCmdPathByName(name string, group string) string {
	return file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "cmd", group, name)
}

func GetScriptAppPathByName(name string) string {
	return file.JoinPath(dir.CoreUtilsShellScripts(), name)
}
