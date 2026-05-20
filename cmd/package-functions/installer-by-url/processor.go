package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/git"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/netc"
	"golangutils/pkg/str"

	"main/cmd/package-functions/installer-by-url/args"
	"main/cmd/package-functions/installer-by-url/entities"
	cmdenums "main/cmd/package-functions/installer-by-url/enums"

	"github.com/spf13/cobra"
)

func github(cmd *cobra.Command) {
	githubArgs := args.GithubArgs{BaseArgs: args.BaseArgs{}}
	githubArgs.FillGithubValues(cmd)
	githubArgs.ValidateGithub()
	data := entities.NewAppInfo(&githubArgs.BaseArgs, destScriptDir, cmdenums.GITHUB_SOURCE)
	data.Github = entities.GithubInfo{
		Owner:      githubArgs.Owner,
		Repository: githubArgs.Repository,
		IsLatest:   githubArgs.IsLatest,
	}
	saveOrCreateConfigFile(*data)
	saveShellScript(githubArgs.ScriptOrCommand, data.FileRunner.ShellScript)
	if githubArgs.RunAfterInit {
		run(data.Name)
	}
}

func runGithub(data *entities.AppInfo) {
	canProcess := false
	url := data.Url
	if urlHasAppVersionKey(data.Url) {
		appVersion, err := git.GithubGetLatestVersionRepo(data.Github.Owner, data.Github.Repository, data.Github.IsLatest, false)
		if err != nil {
			logger.Error(err)
		} else if str.IsEmpty(appVersion.Version) {
			logger.ErrorStr("No releases found for this repository.")
		} else {
			url = getUrlFormated(data.Url, appVersion.Version)
			if appVersion.Version != data.Github.InstalledVersion {
				data.Github.InstalledVersion = appVersion.Version
				canProcess = true
				saveOrCreateConfigFile(*data)
			} else {
				if forceInstallArg {
					canProcess = true
				}
			}
		}
	}
	if canProcess {
		err := netc.Download(url, data.FileRunner.OutputFile)
		if err != nil {
			logger.Error(err)
		} else {
			execInstaller(*data)
		}
	}
}

func generic(cmd *cobra.Command) {
	genericArgs := args.GenericArgs{BaseArgs: args.BaseArgs{}}
	genericArgs.FillGenericValues(cmd)
	genericArgs.ValidateGeneric()
	data := entities.NewAppInfo(&genericArgs.BaseArgs, destScriptDir, cmdenums.GENERIC_SOURCE)
	saveOrCreateConfigFile(*data)
	saveShellScript(genericArgs.ScriptOrCommand, data.FileRunner.ShellScript)
	if genericArgs.RunAfterInit {
		run(data.Name)
	}
}

func runGeneric(data *entities.AppInfo) {
	url := data.Url
	if urlHasAppVersionKey(data.Url) {
		url = getUrlFormated(data.Url, "")
	}
	err := netc.Download(url, data.FileRunner.OutputFile)
	if err != nil {
		logger.Error(err)
	} else {
		execInstaller(*data)
	}
}

func listAll() {
	for index, fileConfig := range getAllFilesConfigs() {
		data, err := file.ReadJsonFile[entities.AppInfo](fileConfig)
		if err != nil {
			logger.Error(err)
		} else {
			count := index + 1
			entry := fmt.Sprintf(`%d - %s`, count, data.Name)
			fmt.Println(entry)
		}
	}
}

func run(appNameArg string) {
	for _, fileConfig := range getAllFilesConfigs() {
		data, err := file.ReadJsonFile[entities.AppInfo](fileConfig)
		if err == nil {
			if str.IsEmpty(appNameArg) || data.Name == appNameArg {
				logger.Header(fmt.Sprintf("Install/Update %s", data.Name))
				switch data.Source {
				case cmdenums.GITHUB_SOURCE:
					runGithub(&data)
				case cmdenums.GENERIC_SOURCE:
					runGeneric(&data)
				default:
					logger.Error(fmt.Errorf("Invalid Source: %s", data.Name))
				}
			}
		} else {
			logger.Error(fmt.Errorf("Failed on read the file: %s", fileConfig))
			logger.Error(err)
		}
	}
}

func uninstall(appNameArg string) {
	for _, fileConfig := range getAllFilesConfigs() {
		data, err := file.ReadJsonFile[entities.AppInfo](fileConfig)
		if err != nil {
			logger.Error(err)
		} else {
			if data.Name == appNameArg {
				if !data.FileRunner.IsCommand {
					logic.ProcessError(file.DeleteFile(data.FileRunner.ShellScript))
				}
				logic.ProcessError(file.DeleteFile(fileConfig))
				logger.Ok(fmt.Sprintf(`Uninstalled %s`, data.Name))
			}
		}
	}
}
