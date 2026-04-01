package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"slices"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	appName   = "alias-manager-cu"
	noAllArgs bool
	forceAgr  bool
)

func init() {
	loadAndValidateVars()
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   appName,
		Short: "Add aliases",
		Long:  "Will be install alias as a script. To get the list of scripts, please run: script-manager-cu list",
	}
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "Name of alias (required)")
	cobralib.CobraCmd.Flags().StringVarP(&contentArg, "command", "c", "", "Command for alias. All args will be passed automatic (required)")
	cobralib.CobraCmd.Flags().BoolVarP(&noAllArgs, "no-args", "o", false, "No insert shell all args(Ex: $args in Powershell)")
	cobralib.CobraCmd.Flags().BoolVarP(&forceAgr, "force", "f", false, "Force install")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("name"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("command"))
	cobralib.WithRun(process)
}

func process() {
	validateName(nameArg)
	if str.IsEmpty(contentArg) {
		logic.ProcessError(fmt.Errorf("invalid content"))
	}
	aliasFile := file.JoinPath(system.TempDir(), file.FileName(nameArg))
	if platform.IsWindows() {
		aliasFile = aliasFile + ".ps1"
	}
	logic.ProcessError(file.DeleteFile(aliasFile))
	fileInfo := models.FileWriterConfig{
		File:        aliasFile,
		IsAppend:    false,
		WithUtf8BOM: false,
		IsCreateDir: true,
	}
	logic.ProcessError(file.Touch(fileInfo.File))
	if !platform.IsWindows() && slices.Contains([]enums.ShellType{enums.Bash, enums.Zsh, enums.Ksh, enums.Fish}, currentShell) {
		fileInfo.Data = fmt.Sprintf("#!/usr/bin/env %s%s%s", currentShell.String(), common.Eol(), common.Eol())
	}
	shellArgs := logic.Ternary(noAllArgs, "", shell.GetShellAllArgsVarStr())
	fileInfo.Data = fmt.Sprintf("%s%s", fileInfo.Data, fmt.Sprintf("%s %s", contentArg, shellArgs))
	logic.ProcessError(file.WriteFile(fileInfo))
	scriptManagerArgs := []string{"install", aliasFile}
	if forceAgr {
		scriptManagerArgs = append(scriptManagerArgs, "-f")
	}
	libs.RunCoreUtilsCmd("script-manager-cu", false, scriptManagerArgs...)
}

func main() {
	cobralib.Run()
}
