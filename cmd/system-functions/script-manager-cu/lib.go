package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"regexp"
	"strings"

	"main/internal/dir"
	"main/internal/libs"
)

func validateOS() {
	if !platform.IsWindows() && !platform.IsLinux() {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
}

func loadAndValidateVars() {
	// Validate
	validateOS()
	// Load others
	config = libs.NewConfig()
	shellDir = dir.CoreUtilsShellScripts()
}

func validateName(name string) {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if str.IsEmpty(name) || !re.MatchString(name) {
		logic.ProcessError(fmt.Errorf("%s: Invalid given name: %s. Name must not be empty and accept only A-Z, a-z, 0-9, - and _", appName, name))
	}
}

func getScriptList() []string {
	files, err := file.ReadDirRecursive(shellDir)
	if err != nil {
		logger.Error(fmt.Errorf("On read directory: %s", shellDir))
		logic.ProcessError(err)
	}
	return files.Files
}

func getScriptVersion(scriptName string) string {
	directory := file.JoinPath(dir.CoreUtilsUserConfig(), "version")
	logic.ProcessError(file.CreateDirectory(directory, true))
	if config.ExistScriptPackageManager(scriptName) {
		versionFilepath := file.JoinPath(directory, scriptName)
		if file.IsFile(versionFilepath) {
			version, _ := file.ReadFile(versionFilepath)
			return strings.Trim(version, common.Eol())
		}
	}
	return ""
}
