package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"regexp"

	"main/internal/dir"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	nameArg    string
	isEnable   bool
	scriptsDir string

	appName = "disable-system-alias-cu"
)

func init() {
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   appName,
		Short: "Manage system aliases",
	}
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "Name of alias (required)")
	cobralib.CobraCmd.Flags().BoolVarP(&isEnable, "is-enable", "e", false, "If system aliases was disabled by CoreUtils, this flag will rollback.")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("name"))
	cobralib.WithRun(process)
}

func validateName(name string) {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if str.IsEmpty(name) || !re.MatchString(name) {
		logic.ProcessError(fmt.Errorf("%s: Invalid given name: %s. Name must not be empty and accept only A-Z, a-z, 0-9, - and _", appName, name))
	}
}

func getScriptFilePath(filename string) string {
	scriptsDir = file.JoinPath(dir.CoreUtilsLocal(), "system-aliases")
	logic.ProcessError(file.CreateDirectory(scriptsDir, true))
	filePath := file.JoinPath(scriptsDir, fmt.Sprintf("system-alias-%s", filename))
	if !file.IsFile(filePath) {
		logic.ProcessError(file.Touch(filePath))
		if !platform.IsWindows() {
			logic.ProcessError(exe.Chmod777(filePath, false))
		}
	}
	return filePath
}

func writeToScriptFile(filePath string, data string) {
	fileInfo := models.FileWriterConfig{
		File:        filePath,
		Data:        data,
		IsAppend:    true,
		WithUtf8BOM: false,
		IsCreateDir: true,
	}
	logic.ProcessError(file.WriteFile(fileInfo))
}

func enableAlias(filePath, data string) {
	if file.IsFile(filePath) {
		logger.Error(file.DeleteFileLines(filePath, data, false))
	}
}

func processPwsh(isDelete bool) {
	powershellFile := getScriptFilePath("pwsh.ps1")
	powershellAliasData := fmt.Sprintf(`if (Get-Item "Alias:\%s" -ErrorAction SilentlyContinue) { Remove-Item "Alias:\%s" -Force }`, nameArg, nameArg)
	powershellFuncData := fmt.Sprintf(`if (Get-Item "Function:\%s" -ErrorAction SilentlyContinue) { Remove-Item "Function:\%s" -Force }`, nameArg, nameArg)
	if isDelete {
		enableAlias(powershellFile, powershellAliasData)
		enableAlias(powershellFile, powershellFuncData)
	} else {
		writeToScriptFile(powershellFile, powershellAliasData)
		writeToScriptFile(powershellFile, powershellFuncData)
	}
}

func process() {
	validateName(nameArg)
	bashData := fmt.Sprintf(`[[ $(alias %s 2>/dev/null) ]] && unalias %s`, nameArg, nameArg)
	zshData := bashData
	kshData := bashData
	fishData := fmt.Sprintf(`functions -e %s`, nameArg)
	bashFile := getScriptFilePath("bash")
	zshFile := getScriptFilePath("zsh")
	kshFile := getScriptFilePath("ksh")
	fishFile := getScriptFilePath("fish.fish")
	//Enable
	processPwsh(true)
	enableAlias(bashFile, bashData)
	enableAlias(zshFile, zshData)
	enableAlias(kshFile, kshData)
	enableAlias(fishFile, fishData)
	//Disable
	if !isEnable {
		processPwsh(false)
		writeToScriptFile(bashFile, bashData)
		writeToScriptFile(zshFile, zshData)
		writeToScriptFile(kshFile, kshData)
		writeToScriptFile(fishFile, fishData)
		logger.Ok(fmt.Sprintf(`Disabled system alias: %s`, nameArg))
	} else {
		logger.Ok(fmt.Sprintf(`Enabled system alias(rollback): %s`, nameArg))
	}
}

func main() {
	cobralib.Run()
}
