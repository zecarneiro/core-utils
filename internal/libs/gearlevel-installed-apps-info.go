package libs

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/obj"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"slices"
)

const GEAR_LEVER_NOT_REGISTED_MSG = "(Not registed by Gear Lever, but the AppImage is installed)"

type GearLeverInstalledAppInfo struct {
	Name          string `json:"name"`
	AppimageFile  string `json:"path"`
	Id            string `json:"desktop_id"`
	Version       string `json:"current_version"`
	IsNotRegisted bool
}

func (i *GearLeverInstalledAppInfo) GetDesktopFile() string {
	fullDesktoFile := fmt.Sprintf("%s/share/applications/%s", system.HomeUserLocalDir(), i.Id)
	if file.IsFile(fullDesktoFile) {
		return fullDesktoFile
	}
	return ""
}

func (i *GearLeverInstalledAppInfo) GetRegisterName() string {
	if !str.IsEmpty(i.Id) {
		return file.FileName(i.Id)
	} else if !str.IsEmpty(i.AppimageFile) {
		return file.FileName(i.AppimageFile)
	}
	return ""
}

type GearLeverInstalledAppsInfo struct {
	SchemaVersion int                         `json:"schema_version"`
	Installed     []GearLeverInstalledAppInfo `json:"installed"`
}

func updateGearLeverListOfInstalledAppsByAppImageFiles(installedListData *GearLeverInstalledAppsInfo) {
	var newApp GearLeverInstalledAppInfo
	appImageDir := fmt.Sprintf(`%s/AppImages`, system.HomeDir())
	dirInfo, err := file.ReadDir(appImageDir)
	isFileRegisted := func(installed GearLeverInstalledAppInfo) bool {
		return installed.AppimageFile == newApp.AppimageFile
	}
	if err == nil {
		for _, fileDir := range dirInfo.Files {
			fileName := file.FileName(fileDir)
			newApp = GearLeverInstalledAppInfo{
				Name:          fmt.Sprintf(`%s %s`, fileName, GEAR_LEVER_NOT_REGISTED_MSG),
				AppimageFile:  fmt.Sprintf(`%s/%s`, appImageDir, fileDir),
				Id:            fmt.Sprintf(`%s.desktop`, fileName),
				IsNotRegisted: true,
			}
			if !slices.ContainsFunc(installedListData.Installed, isFileRegisted) {
				installedListData.Installed = append(installedListData.Installed, newApp)
			}
		}
	}
}

func GetGearLeverListOfInstalledApps() GearLeverInstalledAppsInfo {
	cmdInfo := models.Command{Cmd: "flatpak run it.mijorus.gearlever --list-installed --json"}
	cmdRes, err := exe.Exec(cmdInfo)
	logic.ProcessError(err)
	data, err := obj.StringToObject[GearLeverInstalledAppsInfo](cmdRes)
	logic.ProcessError(err)
	updateGearLeverListOfInstalledAppsByAppImageFiles(&data)
	return data
}
