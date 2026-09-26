package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"
	"main/internal/libs"
)

type AppManifest struct {
	Bin       any        `json:"bin,omitempty"` // string, []string or [][]string
	Shortcuts [][]string `json:"shortcuts,omitempty"`
}

func getManifestFilePath(appRootDir string) string {
	manifestFileName := "manifest.json"
	manifestFile := file.JoinPath(appRootDir, manifestFileName)
	if !file.IsFile(manifestFile) {
		manifestFile = file.JoinPath(appRootDir, fmt.Sprintf(`scoop-%s`, manifestFileName))
	}
	return manifestFile
}

func getScoopAppInfo(appName string) (string, string, string) {
	var exeFileName, exePath string
	appRootDir := libs.RunCoreUtilsCmdWithOutput("scoop-get-app-dir", false, appName)
	manifestData, err := file.ReadJsonFile[AppManifest](getManifestFilePath(appRootDir))
	logic.ProcessError(err)
	if len(manifestData.Shortcuts) > 0 {
		for _, shortcuts := range manifestData.Shortcuts {
			if len(shortcuts) > 0 {
				exeFileName = shortcuts[0]
			}
		}
	}
	if str.IsEmpty(exeFileName) {
		switch bin := manifestData.Bin.(type) {
		case string:
			exeFileName = bin
		case []string:
			exeFileName = logic.Ternary(len(bin) > 0, bin[0], "")
		default:
			exeFileName = ""
		}
	}
	if str.IsEmpty(exeFileName) {
		logic.ProcessError(fmt.Errorf("Can not found executable file name from manifest.json"))
	}
	exePath = file.JoinPath(appRootDir, exeFileName)
	if !file.IsFile(exePath) {
		logic.ProcessError(fmt.Errorf("Can not found executable from manifest.json"))
	}
	return appRootDir, exeFileName, exePath
}
