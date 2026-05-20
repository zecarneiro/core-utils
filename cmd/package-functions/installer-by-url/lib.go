package main

import (
	"fmt"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"strings"

	"main/cmd/package-functions/installer-by-url/entities"
	cmdenums "main/cmd/package-functions/installer-by-url/enums"
)

func getConfigFile(appName string) string {
	return file.JoinPath(configFileDir, fmt.Sprintf("%s.json", appName))
}

func saveShellScript(srcShellScript string, destShellScript string) {
	if file.IsFile(srcShellScript) {
		logic.ProcessError(file.CopyFile(srcShellScript, destShellScript))
		if platform.IsLinux() {
			exe.Chmod777(destShellScript, false)
		}
	}
}

func getUrlFormated(url string, versionApp string) string {
	return str.StringReplaceAll(url, map[string]string{cmdenums.APP_VERSION_KEY: versionApp})
}

func saveOrCreateConfigFile(appInfo entities.AppInfo) {
	logic.ProcessError(file.WriteJsonFile(getConfigFile(appInfo.Name), appInfo, false))
}

func urlHasAppVersionKey(url string) bool {
	return strings.Contains(url, cmdenums.APP_VERSION_KEY)
}

func execInstaller(appInfo entities.AppInfo) {
	shellType := enums.GetShellTypeFromValue(appInfo.FileRunner.Shell)
	exeInfo := models.Command{UseShell: true, Verbose: true}
	if shellType != enums.UnknownShell {
		exeInfo.ShellToUse = shellType
	}
	if file.IsFile(appInfo.FileRunner.OutputFile) {
		if appInfo.FileRunner.IsCommand {
			exeInfo.Cmd = str.StringReplaceAll(appInfo.FileRunner.Command, map[string]string{cmdenums.DOWNLOADED_FILE_ARG_KEY: appInfo.FileRunner.OutputFile})
		} else {
			exeInfo.Cmd = fmt.Sprintf(`%s "%s"`, appInfo.FileRunner.ShellScript, appInfo.FileRunner.OutputFile)
		}
		logger.Error(exe.ExecRealTime(exeInfo))
	} else {
		logger.Error(fmt.Errorf("Downloaded file '%s' not found.", appInfo.FileRunner.OutputFile))
	}
}

func getAllFilesConfigs() []string {
	files := []string{}
	filesInfo, err := file.ReadDir(configFileDir)
	logic.ProcessError(err)
	for _, info := range filesInfo.Files {
		files = append(files, file.JoinPath(configFileDir, info))
	}
	return files
}
