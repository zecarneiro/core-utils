package processors

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"slices"
	"strings"

	"main/cmd/tools-functions/vscode-config/common"
	"main/cmd/tools-functions/vscode-config/entities"
)

type Profiles struct {
	profiles []entities.Profile
}

func NewProfileProcessor() *Profiles {
	processor := &Profiles{}
	processor.loadData()
	return processor
}

func (pp *Profiles) loadData() {
	if common.IsJsonFileInstalled() {
		pp.profiles = common.Configurations.Profiles
	}
}

func (pp *Profiles) ProfileExists(name string) bool {
	command := common.GetVscodeListExtCommand(name)
	command.Verbose = false
	result, _ := exe.Exec(command)
	return logic.Ternary(strings.Contains(result, "Profile '"+name+"' not found."), false, true)
}

func (pp *Profiles) getAllProfileData(name string) []entities.ProfileData {
	profileData := []entities.ProfileData{}
	for _, profile := range common.Configurations.Profiles {
		if profile.Name == name {
			profileData = append(profile.Extensions, pp.processDependsProfiles(profile.Name, profile.DependsProfile, false)...)
			break
		}
	}
	return profileData
}

func (pp *Profiles) processDependsProfiles(profileName string, dependsName []string, verbose bool) []entities.ProfileData {
	listExtensions := []entities.ProfileData{}
	if len(dependsName) > 0 {
		for _, profile := range pp.profiles {
			if profile.Name != profileName && pp.isValidProfileName(profile.Name) && slices.Contains(dependsName, profile.Name) {
				if verbose {
					logger.Header("Append extensions from Profile: " + profile.Name)
				}
				listExtensions = append(listExtensions, profile.Extensions...)
			}
		}
	}
	return listExtensions
}

func (pp *Profiles) isValidProfileName(name string) bool {
	isValid := true
	if strings.Contains(name, " ") {
		isValid = false
		logger.ErrorStr("Invalid profile name: " + name)
		logger.Info("Name must not contains: space")
	}
	return isValid
}

func (pp *Profiles) createProfile(name string) {
	if !pp.ProfileExists(name) {
		var commandStr string
		if platform.IsLinux() {
			commandStr = "(code --profile \"%s\" &) && sleep 2 && killall -9 code && sleep 2"
		} else if platform.IsWindows() {
			commandStr = "Start-Process code -ArgumentList '--profile \"%s\"'; Start-Sleep -Seconds 2; Get-Process code -ErrorAction SilentlyContinue | Stop-Process -Force; Start-Sleep -Seconds 2"
		} else {
			commandStr = ""
			logger.ErrorStr("Can not create this profile, because this SO is not supported")
			logger.Info(fmt.Sprintf("Please open VSCode and create this profile: %s", name))
			console.Pause()
		}
		if len(commandStr) > 0 {
			logger.Header("Creating Profile")
			command := models.Command{Cmd: fmt.Sprintf(commandStr, name), UseShell: true, Verbose: false}
			logic.ProcessError(exe.ExecRealTime(command))
		}
	}
}

func (pp *Profiles) GetAllInstallSettings(profileName string, verbose bool) map[string]interface{} {
	settings := make(map[string]interface{})
	profileSelectedList := slice.FilterArray(common.Configurations.Profiles, func(profile entities.Profile) bool {
		if str.IsEmpty(profileName) {
			return logic.Ternary(profile.IsSettingName, true, false)
		}
		return logic.Ternary(profile.Name == profileName, true, false)
	})
	profileSelected := logic.Ternary(len(profileSelectedList) > 0, profileSelectedList[0], entities.Profile{Name: "", DependsProfile: []string{}})
	for _, profile := range common.Configurations.Profiles {
		if slices.Contains(profileSelected.DependsProfile, profile.Name) {
			if verbose {
				logger.Debug(fmt.Sprintf("Import settings from profile: %s", profile.Name))
			}
			settings = slice.ConcatMap(settings, profile.Settings)
		} else if profile.Name == profileSelected.Name {
			if verbose {
				logger.Debug(fmt.Sprintf("Import self settings: %s", profile.Name))
			}
			settings = slice.ConcatMap(settings, profile.Settings)
		}
	}
	return settings
}

func (pp *Profiles) GetAllProfile() []entities.ProfileStatus {
	status := []entities.ProfileStatus{}
	for _, profile := range common.Configurations.Profiles {
		name := profile.Name
		if profile.IsSettingName {
			name = common.Configurations.SettingsName
		}
		status = append(status, entities.ProfileStatus{Name: name, IsInstalled: pp.ProfileExists(name)})
	}
	return status
}

func (pp *Profiles) GetAllExtensionsFromProfile(name string) []string {
	extensions := []string{}
	for _, data := range pp.getAllProfileData(name) {
		extensions = append(extensions, data.Ids...)
	}
	return extensions
}
