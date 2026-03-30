package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"golangutils/pkg/env"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"strings"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	nameArg        string
	execArg        string
	argsArg        []string
	iconArg        string
	categoriesArg  string
	commentArg     string
	terminalArg    bool
	windowStyleArg string
	runAsAdminArg  bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "create-menu-entry",
		Short: "Create a menu entry (.desktop file on Linux)",
		Long:  "Creates a menu entry in the system",
	}
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "Application name (required)")
	cobralib.CobraCmd.Flags().StringVarP(&execArg, "exec", "e", "", "Executable path (required)")
	cobralib.CobraCmd.Flags().StringArrayVarP(&argsArg, "args", "a", []string{}, "Arguments to pass to the executable")
	cobralib.CobraCmd.Flags().StringVarP(&iconArg, "icon", "i", "", "Icon path")
	cobralib.CobraCmd.Flags().StringVarP(&categoriesArg, "categories", "c", "", "Categories (e.g., Utility;Development)")
	cobralib.CobraCmd.Flags().StringVarP(&commentArg, "comment", "m", "", "Comment/description")
	cobralib.CobraCmd.Flags().BoolVarP(&terminalArg, "terminal", "t", false, "Run in terminal")
	cobralib.CobraCmd.Flags().BoolVarP(&runAsAdminArg, "admin", "A", false, "Run as administrator (Windows SO only)")
	cobralib.WithRun(process)
}

func process() {
	if str.IsEmpty(nameArg) {
		logic.ProcessError(fmt.Errorf("Name is required"))
	}
	if str.IsEmpty(execArg) {
		logic.ProcessError(fmt.Errorf("Exec is required"))
	}
	if platform.IsWindows() {
		createWindowsShortcut()
	} else if platform.IsLinux() {
		processForLinux()
	} else {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
}

func getNameNormalized() string {
	return strings.ToLower(strings.ReplaceAll(nameArg, " ", "-"))
}

func processForLinux() {
	desktopDir := file.JoinPath(system.HomeDir(), ".local/share/applications")
	filename := getNameNormalized() + ".desktop"
	isTerminal := logic.Ternary(terminalArg, "true", "false")
	entryFilepath := file.JoinPath(desktopDir, filename)
	logic.ProcessError(file.CreateDirectory(desktopDir, true))
	fileConfig := models.FileWriterConfig{
		File:        entryFilepath,
		Data:        fmt.Sprintf(linuxTermplate, nameArg, execArg, slice.ArrayToString(argsArg), iconArg, categoriesArg, commentArg, isTerminal),
		IsAppend:    false,
		WithUtf8BOM: false,
	}
	logic.ProcessError(file.WriteFile(fileConfig))
	logger.Ok(fmt.Sprintf("Created: %s", entryFilepath))
}

func createWindowsShortcut() {
	eol := common.Eol()
	othersData := ""
	nameNormalized := getNameNormalized()
	envAppData := env.Get("APPDATA")
	if len(envAppData) == 0 {
		logic.ProcessError(fmt.Errorf("Not found APPDATA env. Please, verify your system."))
	}
	startMenuDir := fmt.Sprintf("%s\\Microsoft\\Windows\\Start Menu\\Programs", envAppData[0])
	filename := nameNormalized + ".lnk"
	entryFilepath := file.JoinPath(startMenuDir, filename)

	targetPath := execArg
	arguments := slice.ArrayToString(argsArg)
	windowStyle := 1
	switch strings.ToLower(windowStyleArg) {
	case "maximized":
		windowStyle = 3
	case "minimized":
		windowStyle = 7
	}
	if terminalArg {
		targetPath = "cmd.exe"
		arguments = fmt.Sprintf("/c `\"%s %s`\"", execArg, slice.ArrayToString(argsArg))
	}
	if !str.IsEmpty(commentArg) {
		othersData = fmt.Sprintf("%s%s$Shortcut.Description = \"%s\"", othersData, eol, commentArg)
	}
	if !str.IsEmpty(iconArg) {
		othersData = fmt.Sprintf("%s%s$Shortcut.IconLocation = \"%s\"", othersData, eol, iconArg)
	}
	if !str.IsEmpty(arguments) {
		othersData = fmt.Sprintf("%s%s$Shortcut.Arguments = \"%s\"", othersData, eol, arguments)
	}
	scriptData := fmt.Sprintf(windowsTemplate, entryFilepath, targetPath, windowStyle, othersData)
	scriptPath := file.JoinPath(system.TempDir(), fmt.Sprintf("create_shortcut-%s.ps1", nameNormalized))
	fileConfig := models.FileWriterConfig{
		File:        scriptPath,
		Data:        scriptData,
		IsAppend:    false,
		WithUtf8BOM: false,
	}
	logic.ProcessError(file.WriteFile(fileConfig))
	libs.RunCoreUtilsCmd("run-shell-script", false, "-f", scriptPath, "-s", enums.PowerShell.String())
	logic.ProcessError(file.DeleteFile(scriptPath))
	logger.Ok(fmt.Sprintf("Created shortcut: %s", entryFilepath))
	if runAsAdminArg {
		adminScriptData := fmt.Sprintf(windowsAdminTemplate, entryFilepath)
		adminScriptPath := file.JoinPath(system.TempDir(), fmt.Sprintf("set_admin-%s.ps1", nameNormalized))
		fileConfigAdmin := models.FileWriterConfig{
			File:        adminScriptPath,
			Data:        adminScriptData,
			IsAppend:    false,
			WithUtf8BOM: false,
		}
		logic.ProcessError(file.WriteFile(fileConfigAdmin))
		libs.RunCoreUtilsCmd("run-shell-script", false, "-f", adminScriptPath, "-s", enums.PowerShell.String())
		logic.ProcessError(file.DeleteFile(adminScriptPath))
		logger.Ok(fmt.Sprintf("Set as administrator: %s", entryFilepath))
	}
}

func main() {
	cobralib.Run()
}
