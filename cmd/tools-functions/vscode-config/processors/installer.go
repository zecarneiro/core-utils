package processors

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/obj"
	"slices"
	"strconv"
	"strings"

	"main/cmd/tools-functions/vscode-config/common"
	"main/cmd/tools-functions/vscode-config/entities"
	"main/internal/libs"
)

type Installer struct {
	profilesProcessor     *Profiles
	downloaderProcessor   *Downloader
	profileAlreadyCreated []string
}

func NewInstallProcessor(profileProcessor *Profiles) *Installer {
	return &Installer{
		profilesProcessor:     profileProcessor,
		downloaderProcessor:   &Downloader{},
		profileAlreadyCreated: []string{},
	}
}

/* --------------------------------- PRIVATE -------------------------------- */
const (
	MAX_INSTALL_EXTENSIONS = 5
)

func (ip *Installer) InstallJsonFile(jsonFile string) {
	jsonFile = file.ResolvePath(jsonFile)
	if file.IsFile(jsonFile) {
		logger.Info("Copy given file to: " + common.JsonFile)
		logic.ProcessError(file.CopyFile(jsonFile, common.JsonFile))
	}
}

func (ip *Installer) ProcessInstall() {
	skipProfiles := []string{}
	console.WaitForAnyKeyPressed("Please, close all instance of VSCode and PRESS ANY KEY TO CONTINUE...")
	for _, profile := range ip.profilesProcessor.profiles {
		if !profile.CanInstall {
			skipProfiles = append(skipProfiles, profile.Name)
			continue
		}
		logger.Log("")
		logger.Title("Process Profile: " + profile.Name)
		profile.Extensions = ip.profilesProcessor.getAllProfileData(profile.Name)
		if ip.profilesProcessor.isValidProfileName(profile.Name) {
			ip.processByProfile(profile)
		}
		if profile.CanInstall {
			logger.Header("Processing Profile: " + profile.Name + ", Done.")
		}
	}
	ip.SetSettingConfigurations()
	if len(skipProfiles) > 0 {
		logger.Header("SKIPED PROFILES TO INSTALL")
		for _, profile := range skipProfiles {
			logger.Log("- " + profile)
		}
	}
}

func (ip *Installer) processByProfile(profile entities.Profile) {
	profileName := profile.Name
	ip.downloaderProcessor.force = false
	if profile.IsSettingName {
		profileName = common.Configurations.SettingsName
	}
	ip.profilesProcessor.createProfile(profileName)
	ip.profileAlreadyCreated = append(ip.profileAlreadyCreated, profileName)
	counter := 1
	for {
		extensionsToInstall := ip.getExtensionsToInstall(profileName, profile.Extensions)
		if len(extensionsToInstall) == 0 {
			break
		} else {
			logger.Header("Download Extensions")
			ip.downloaderProcessor.downloadList(extensionsToInstall)
		}
		if counter > MAX_INSTALL_EXTENSIONS && len(extensionsToInstall) > 0 {
			if ip.askTryInstallAllExtensions(extensionsToInstall) {
				ip.downloaderProcessor.force = true
				if !profile.IsSettingName {
					ip.profilesProcessor.createProfile(profileName)
				}
				counter = 1
			} else {
				break
			}
		}
		logger.Header("Install Extensions")
		for _, data := range profile.Extensions {
			for _, id := range data.Ids {
				if slices.Contains(extensionsToInstall, id) {
					ip.installExtension(profileName, id)
				}
			}
		}
		counter++
	}
}

func (ip *Installer) SetSettingConfigurations() {
	settingsDir := libs.GetVscodeSettingsDir()
	settingsFileName := "settings.json"
	logger.Separator()
	logger.Header("Set settings")
	settings := ip.profilesProcessor.GetAllInstallSettings("", true)
	if len(settingsDir) > 0 {
		logic.ProcessError(file.CreateDirectory(settingsDir, true))
		logic.ProcessError(file.WriteJsonFile(file.JoinPath(settingsDir, settingsFileName), settings, true))
	} else {
		logger.Debug("\n\nGo to setting and open json settings")
		logger.Debug("Append this setting bellow on json settings")
		settingsStr, err := obj.ObjectToString(settings)
		logic.ProcessError(err)
		fmt.Println(settingsStr)
		common.OpenVscodeWithNewProfile(common.Configurations.SettingsName)
		console.Pause()
	}
	logger.Separator()
	logger.Header("Apply all settings for all Profiles")
	for _, userDataProfile := range libs.GetVscodeUserDataProfiles() {
		logger.Info(fmt.Sprintf("Set settings for: PROFILE_NAME(%s) and PROFILE_ID(%s)", userDataProfile.Name, userDataProfile.Location))
		settings = ip.profilesProcessor.GetAllInstallSettings(userDataProfile.Name, true)
		profileDir := file.JoinPath(settingsDir, "profiles", userDataProfile.Location)
		logic.ProcessError(file.WriteJsonFile(file.JoinPath(profileDir, settingsFileName), settings, true))
		logger.Log("")
	}
}

func (ip *Installer) getExtensionsToInstall(profileName string, extensions []entities.ProfileData) []string {
	listExtensions := ip.getInstalledExtensions(profileName)
	listToInstal := []string{}
	for _, data := range extensions {
		for _, id := range data.Ids {
			if !ip.isExtensionInstalled(id, listExtensions) {
				listToInstal = append(listToInstal, id)
			}
		}
	}
	return listToInstal
}

func (ip *Installer) askTryInstallAllExtensions(extensions []string) bool {
	fmt.Println("\n####### NOT INSTALLED EXTENSIONS ID'S #######")
	for index, id := range extensions {
		fmt.Println(strconv.Itoa(index) + " - " + id)
	}
	return console.Confirm("Continue", false)
}

func (ip *Installer) installExtension(profileName string, id string) {
	command := common.GetCodeCommand()
	fileId := ip.downloaderProcessor.getExtensionVsixFile(id)
	if !file.IsFile(fileId) {
		ip.downloaderProcessor.download(id)
		if file.IsFile(fileId) {
			id = fileId
		}
	} else {
		id = fileId
	}
	command.Args = append(command.Args, profileName, "--install-extension", id)
	command.Verbose = true
	resp, err := exe.Exec(command)
	logger.Error(err)
	if !strings.Contains(resp, "was successfully installed") {
		logger.ErrorStr(resp)
	}
}

func (ip *Installer) getInstalledExtensions(profileName string) []string {
	listExtensions := []string{}
	command := common.GetVscodeListExtCommand(profileName)
	output, err := exe.Exec(command)
	if err != nil {
		logger.Error(err)
		return listExtensions
	}
	if len(output) > 0 {
		for _, extension := range strings.Split(output, "\n") {
			listExtensions = append(listExtensions, strings.ToLower(strings.TrimSpace(extension)))
		}
	}
	return listExtensions
}

func (ip *Installer) isExtensionInstalled(id string, extensionsInstalled []string) bool {
	return slices.Contains(extensionsInstalled, strings.ToLower(id))
}
