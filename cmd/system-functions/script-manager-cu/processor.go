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
		forceArg, errForce := cmd.Flags().GetBool("force")
		logic.ProcessError(errForce)
		dest := file.JoinPath(shellDir, logic.Ternary(platform.IsWindows(), file.Basename(path), filename))
		canInstall := true
		if !forceArg && file.IsFile(dest) {
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
	if !str.IsEmpty(name) {
		validateName(name)
	}
	// Process
	logger.Title("Run all script(s) to install/update/remove package(s)")
	for _, script := range getScriptList() {
		canRun := false
		scriptFileName := file.FileName(script)
		if config.ExistScriptPackageManager(scriptFileName) {
			canRun = true
			basename := file.Basename(script)
			if !str.IsEmpty(name) && !(name == script || name == basename || name == scriptFileName) {
				canRun = false
			}
		}
		if canRun {
			logger.Header(fmt.Sprintf("Running %s", scriptFileName))
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
			packageManagerInfo := logic.Ternary(config.ExistScriptPackageManager(file.Basename(script)), "(Package Manager)", "")
			logger.Log(fmt.Sprintf("%d. %s %s", count, file.FileName(script), packageManagerInfo))
			count++
		}
	}
}

func update() {
	if platform.IsWindows() {
		for _, script := range getScriptList() {
			libs.CreateExecPwshFromPromptCMD(script)
		}
	}
}

func scriptVersionManager(cmd *cobra.Command, args []string) {
	scriptArg := file.FileName(slice.ArrayToString(args))
	versionArg, err := cmd.Flags().GetString("version")
	logic.ProcessError(err)
	directory := file.JoinPath(dir.CoreUtilsUserConfig(), "version")
	logic.ProcessError(file.CreateDirectory(directory, true))
	if config.ExistScriptPackageManager(scriptArg) {
		versionFilepath := file.JoinPath(directory, scriptArg)
		scriptVersion := ""
		if file.IsFile(versionFilepath) {
			version, _ := file.ReadFile(versionFilepath)
			scriptVersion = version
		}
		if str.IsEmpty(versionArg) {
			fmt.Println(scriptVersion)
		} else {
			if versionArg != versionFilepath {
				file.WriteFile(models.FileWriterConfig{File: versionFilepath, Data: versionArg, IsAppend: false, IsCreateDir: true, WithUtf8BOM: false})
			}
		}
	}
}
