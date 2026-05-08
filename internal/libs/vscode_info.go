package libs

import (
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
)

type VSCodeUserDataProfiles struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type vscodeStorageInfo struct {
	UserDataProfiles []VSCodeUserDataProfiles `json:"userDataProfiles"`
}

func GetVscodeSettingsDir() string {
	settingsDir := ""
	if platform.IsWindows() {
		settingsDir = file.JoinPath(system.HomeDir(), "AppData\\Roaming\\Code\\User")
	} else if platform.IsLinux() {
		settingsDir = file.JoinPath(system.HomeDir(), ".config/Code/User")
	}
	return settingsDir
}

func GetVscodeUserDataProfiles() []VSCodeUserDataProfiles {
	settingsDir := GetVscodeSettingsDir()
	storageJsonFile := file.JoinPath(settingsDir, "globalStorage", "storage.json")
	if file.IsFile(storageJsonFile) {
		vscodeStorageInfo, err := file.ReadJsonFile[vscodeStorageInfo](file.JoinPath(settingsDir, "globalStorage", "storage.json"))
		logic.ProcessError(err)
		return vscodeStorageInfo.UserDataProfiles
	}
	return []VSCodeUserDataProfiles{}
}
