package main

import (
	"fmt"
	"golangutils/pkg/conv"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/obj"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"os"

	"main/cmd/tools-functions/vscode-config/common"
	"main/cmd/tools-functions/vscode-config/entities"
	"main/cmd/tools-functions/vscode-config/processors"
)

type Processor struct {
	jsonFileArg           string
	settingsArg           bool
	processAllArg         bool
	extractSettingsArg    bool
	extractExtsArg        bool
	extractAllArg         bool
	extractProfileNameArg string

	installerProcessor *processors.Installer
	profilesProfessor  *processors.Profiles
}

func NewProcessor() *Processor {
	p := &Processor{}
	p.loadData()
	return p
}

/* ------------------------------ PRIVATE AREA ------------------------------ */
func (p *Processor) loadData() {
	common.JsonInfo = &entities.JsonInfo{}
	common.FillJsonFile(false)
	p.profilesProfessor = processors.NewProfileProcessor()
	p.installerProcessor = processors.NewInstallProcessor(p.profilesProfessor)
}

func (p *Processor) resetVscode() {
	var pathsCmd []string
	if platform.IsLinux() {
		pathsCmd = []string{
			file.JoinPath(system.HomeUserConfigDir(), "Code"),
			file.JoinPath(system.HomeDir(), ".vscode"),
			file.JoinPath(system.HomeDir(), ".cache/code"),
		}
	} else if platform.IsWindows() {
		pathsCmd = []string{
			file.JoinPath(os.Getenv("APPDATA"), "Code"),
			file.JoinPath(system.HomeDir(), ".vscode"),
			file.JoinPath(os.Getenv("LOCALAPPDATA"), "Code"),
		}
	} else if platform.IsDarwin() {
		pathsCmd = []string{
			file.JoinPath(system.HomeDir(), "Library/Application Support/Code"),
			file.JoinPath(system.HomeDir(), ".vscode"),
			file.JoinPath(system.HomeDir(), "Library/Caches/com.microsoft.VSCode"),
			file.JoinPath(system.HomeDir(), "Library/Preferences/com.microsoft.VSCode.plist"),
		}
	} else {
		logger.ErrorStr(platform.InvalidMSG)
		pathsCmd = []string{}
	}
	for _, pathCmd := range pathsCmd {
		if file.IsFile(pathCmd) {
			logger.Info(fmt.Sprintf("Deleting file: %s", pathCmd))
			logger.Error(file.DeleteFile(pathCmd))
		}
		if file.IsDir(pathCmd) {
			logger.Info(fmt.Sprintf("Deleting directory: %s", pathCmd))
			logger.Error(file.DeleteDirectory(pathCmd))
		}
	}
}

func (p *Processor) listProfiles() {
	for _, profile := range p.profilesProfessor.GetAllProfile() {
		if profile.IsInstalled {
			logger.Log("- " + profile.Name)
		} else {
			logger.Log("- " + profile.Name + " (Not Installed)")
		}
	}
}

func (p *Processor) profileStatus(name string) {
	fmt.Println(conv.BoolToString(processor.profilesProfessor.ProfileExists(name)))
}

func (p *Processor) install() {
	if !str.IsEmpty(p.jsonFileArg) {
		p.installerProcessor.InstallJsonFile(p.jsonFileArg)
	}
	common.FillJsonFile(true)
	p.loadData()
	if p.processAllArg {
		p.installerProcessor.ProcessInstall()
	} else if p.settingsArg {
		p.installerProcessor.SetSettingConfigurations()
	}
}

func (p *Processor) extract() {
	dataObj := map[string]any{}
	if str.IsEmpty(p.extractProfileNameArg) {
		logic.ProcessError(fmt.Errorf("Invalid given profile name"))
	}
	if p.extractAllArg || p.extractSettingsArg {
		dataObj["settings"] = p.profilesProfessor.GetAllInstallSettings(p.extractProfileNameArg, false)
	}
	if p.extractAllArg || p.extractExtsArg {
		dataObj["extensions"] = p.profilesProfessor.GetAllExtensionsFromProfile(p.extractProfileNameArg)
	}
	data, err := obj.ObjectToString(dataObj)
	logic.ProcessError(err)
	logger.Log(data)
}
