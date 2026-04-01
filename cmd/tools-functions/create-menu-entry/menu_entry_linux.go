//go:build linux

package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/conv"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
)

var template = `[Desktop Entry]
Version=1.0
Type=Application
Name=%s
Exec=%s
Icon=%s
Categories=%s
Comment=%s
Terminal=%s
X-AppStream-Ignore=true
`

type MenuEntryLinuxProcessor struct {
	desktopDir   string
	desktopFile  string
	templateData string
}

func (m *MenuEntryLinuxProcessor) loadVars() {
	desktopFilename := getNameNormalized(menuEntryData.name, false) + ".desktop"
	m.desktopDir = file.JoinPath(system.HomeDir(), ".local/share/applications")
	m.desktopFile = file.JoinPath(m.desktopDir, desktopFilename)
}

func (m *MenuEntryLinuxProcessor) buildTemplate() {
	// Exec
	exec := logic.Ternary(str.IsEmpty(menuEntryData.argsStr), menuEntryData.exec, fmt.Sprintf(`%s %s`, menuEntryData.exec, menuEntryData.argsStr))
	// Terminal
	terminal := conv.BoolToString(menuEntryData.terminal)
	icon := menuEntryData.icon
	if menuEntryData.terminal && str.IsEmpty(icon) {
		icon = "utilities-terminal"
	}
	// Build
	m.templateData = fmt.Sprintf(template, menuEntryData.name, exec, icon, menuEntryData.categories, menuEntryData.comment, terminal)
}

func start() {
	if !platform.IsLinux() {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
	processor := MenuEntryLinuxProcessor{}
	processor.loadVars()
	logic.ProcessError(file.CreateDirectory(processor.desktopDir, true))
	processor.buildTemplate()
	writeToFile(processor.desktopFile, processor.templateData)
	printOkMsg(processor.desktopFile)
}
