//go:build windows

package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/enums"
	"golangutils/pkg/env"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"main/internal/libs"
)

var (
	scriptTemplate = `
$WshShell = New-Object -ComObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("%s")
$Shortcut.TargetPath = "%s"
$Shortcut.WindowStyle = 1
%s
%s
%s
$Shortcut.Save()
`
	scriptAdminTemplate = `
$path = "%s"
$bytes = [System.IO.File]::ReadAllBytes($path)
$bytes[0x15] = $bytes[0x15] -bor 0x20
[System.IO.File]::WriteAllBytes($path, $bytes)
`
)

type MenuEntryWindowsProcessor struct {
	windowStyle             string
	startMenuDir            string
	lnkFile                 string
	scriptFile              string
	scriptAdminFile         string
	nameNormalized          string
	scriptTemplateData      string
	scriptAdminTemplateData string
}

func (m *MenuEntryWindowsProcessor) getEnvAppData() string {
	envAppData := env.Get("APPDATA")
	if len(envAppData) == 0 {
		logic.ProcessError(fmt.Errorf("Not found APPDATA env. Please, verify your system."))
	}
	return envAppData[0]
}

func (m *MenuEntryWindowsProcessor) loadVars() {
	lnkFilename := menuEntryData.name + ".lnk"
	m.startMenuDir = fmt.Sprintf("%s\\Microsoft\\Windows\\Start Menu\\Programs", m.getEnvAppData())
	m.lnkFile = file.JoinPath(m.startMenuDir, lnkFilename)
	m.nameNormalized = getNameNormalized(menuEntryData.name, true)
	m.scriptFile = file.JoinPath(system.TempDir(), fmt.Sprintf("create-shortcut-%s.ps1", m.nameNormalized))
	m.scriptAdminFile = file.JoinPath(system.TempDir(), fmt.Sprintf("set-admin-%s.ps1", m.nameNormalized))
}

func (m *MenuEntryWindowsProcessor) getTargetArgs() string {
	isTargetArgsEmpty := str.IsEmpty(menuEntryData.argsStr)
	args := ""
	if menuEntryData.terminal {
		exec := logic.Ternary(isTargetArgsEmpty, menuEntryData.exec, fmt.Sprintf(`%s %s`, menuEntryData.exec, menuEntryData.argsStr))
		args = fmt.Sprintf("/c `\"%s`\"", exec)
	} else if !isTargetArgsEmpty {
		args = menuEntryData.argsStr
	}
	return logic.Ternary(str.IsEmpty(args), "", fmt.Sprintf(`$Shortcut.Arguments = "%s"`, args))
}

func (m *MenuEntryWindowsProcessor) buildTemplate() {
	// TargetPath and targetArgs
	targetPath := menuEntryData.exec
	targetArgs := m.getTargetArgs()
	if menuEntryData.terminal {
		targetPath = "cmd.exe"
	}
	// Description
	description := ""
	if !str.IsEmpty(menuEntryData.comment) {
		description = fmt.Sprintf(`$Shortcut.Description = "%s"`, menuEntryData.comment)
	}
	// IconLocation
	iconLocation := ""
	if !str.IsEmpty(menuEntryData.icon) {
		description = fmt.Sprintf(`$Shortcut.IconLocation = "%s"`, menuEntryData.icon)
	}
	// Build
	m.scriptTemplateData = fmt.Sprintf(scriptTemplate, m.lnkFile, targetPath, targetArgs, description, iconLocation)
}

func (m *MenuEntryWindowsProcessor) buildAdminTemplate() {
	// Build
	m.scriptAdminTemplateData = fmt.Sprintf(scriptAdminTemplate, m.lnkFile)
}

func start() {
	if !platform.IsWindows() {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
	processor := MenuEntryWindowsProcessor{}
	processor.loadVars()

	// Create .lnk file
	logic.ProcessError(file.CreateDirectory(processor.startMenuDir, true))
	processor.buildTemplate()
	writeToFile(processor.scriptFile, processor.scriptTemplateData)
	libs.RunCoreUtilsCmd("run-shell-script", false, "-f", processor.scriptFile, "-s", enums.PowerShell.String())
	logic.ProcessError(file.DeleteFile(processor.scriptFile))
	printOkMsg(processor.lnkFile)

	// Set as Admin
	if menuEntryData.runAsAdmin {
		processor.buildAdminTemplate()
		writeToFile(processor.scriptAdminFile, processor.scriptAdminTemplateData)
		libs.RunCoreUtilsCmd("run-shell-script", false, "-f", processor.scriptAdminFile, "-s", enums.PowerShell.String())
		logic.ProcessError(file.DeleteFile(processor.scriptAdminFile))
		logger.Ok(fmt.Sprintf("Set as administrator: %s", processor.lnkFile))
	}
}
