package libs

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/obj"
	"golangutils/pkg/system"
)

type GearlevelInstalledAppInfo struct {
	Name         string `json:"name"`
	AppimageFile string `json:"path"`
	Id           string `json:"desktop_id"`
	Version      string `json:"current_version"`
}

func (i *GearlevelInstalledAppInfo) GetDesktopFile() string {
	fullDesktoFile := fmt.Sprintf("%s/share/applications/%s", system.HomeUserLocalDir(), i.Id)
	if file.IsFile(fullDesktoFile) {
		return fullDesktoFile
	}
	return ""
}

func (i *GearlevelInstalledAppInfo) GetRegisterName() string {
	return file.FileName(i.Id)
}

type GearlevelInstalledAppsInfo struct {
	SchemaVersion int                         `json:"schema_version"`
	Installed     []GearlevelInstalledAppInfo `json:"installed"`
}

func GetGearlevelListOfInstalledApps() GearlevelInstalledAppsInfo {
	cmdInfo := models.Command{Cmd: "flatpak run it.mijorus.gearlever --list-installed --json"}
	cmdRes, err := exe.Exec(cmdInfo)
	logic.ProcessError(err)
	data, err := obj.StringToObject[GearlevelInstalledAppsInfo](cmdRes)
	logic.ProcessError(err)
	return data
}
