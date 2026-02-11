package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"slices"
	"strings"

	"main/internal/dir"
	"main/internal/libs"

	"github.com/spf13/cobra"
)

func install(cmd *cobra.Command, args []string) {
	path := slice.ArrayToString(args)
	filename := file.FileName(path)
	validateName(filename)
	// Get all cmd args
	isPackageManagerArg, err := cmd.Flags().GetBool("is-package-manager")
	logic.ProcessError(err)
	// Process
	if file.IsFile(path) {
		dest := file.JoinPath(shellDir, logic.Ternary(platform.IsWindows(), file.Basename(path), filename))
		canInstall := true
		if file.IsFile(dest) {
			canInstall = console.Confirm(fmt.Sprintf("Script '%s' already exists. Continue?", dest), true)
		}
		if canInstall {
			logic.ProcessError(file.CopyFile(path, dest))
			if !platform.IsWindows() {
				logic.ProcessError(exe.Chmod777(dest, false))
			}
			if isPackageManagerArg {
				config.ScriptPackageManager = append(config.ScriptPackageManager, filename)
				config.Write()
			}
			logger.Ok(fmt.Sprintf("Script installed with sucess: %s", filename))
		}
	} else {
		logic.ProcessError(fmt.Errorf("Invalid given file: %s", path))
	}
	update()
}

func uninstall(name string) {
	// Validate
	validateName(name)
	// Process
	exist_script := false
	for _, script := range getScriptList() {
		filename := file.FileName(script)
		if name == script || name == filename {
			logic.ProcessError(file.DeleteFile(script))
			logger.Ok(fmt.Sprintf("Deleted: %s", filename))
			exist_script = true
		}
	}
	if !exist_script {
		logic.ProcessError(fmt.Errorf("Not found: %s", name))
	} else {
		update()
	}
}

func run(name string) {
	// Validate
	validateName(name)
	// Process
	logger.Title("Run all script(s) to install/update/remove package(s)")
	for _, script := range getScriptList() {
		basename := file.Basename(script)
		canRun := true
		if !str.IsEmpty(name) && !(name == script || name == file.FileName(basename)) {
			canRun = false
		}
		if slices.Contains(config.ScriptPackageManager, basename) && canRun {
			cmd := models.Command{Cmd: fmt.Sprintf(". \"%s\"", script), UseShell: true, Verbose: true}
			logger.Error(exe.ExecRealTime(cmd))
		}
	}
}

func list(filter string) {
	// Process
	logger.Title("List of all script(s) to install/update/remove package(s)")
	count := 1
	for _, script := range getScriptList() {
		filename := file.FileName(script)
		if len(filter) == 0 || strings.Contains(filename, filter) {
			packageManagerInfo := logic.Ternary(slices.Contains(config.ScriptPackageManager, file.Basename(script)), "(Package Manager)", "")
			logger.Log(fmt.Sprintf("%d. %s %s", count, file.FileName(script), packageManagerInfo))
			count++
		}
	}
}

func update() {
	if platform.IsWindows() {
		promptDir := dir.CoreUtilsPrompt()
		for _, script := range getScriptList() {
			libs.CreateExecPwshFromPromptCMD(script)
		}
		logger.Ok(fmt.Sprintf("Update PROMPT dir with all coreutils scripts: %s", promptDir))
	}
}
