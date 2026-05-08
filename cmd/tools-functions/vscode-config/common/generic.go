package common

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/cmd/tools-functions/vscode-config/entities"
)

func FillJsonFile(verbose bool) {
	if IsJsonFileInstalled() {
		if verbose {
			logger.Info("Reading file: " + JsonFile)
		}
		data, err := file.ReadJsonFile[entities.Configurations](JsonFile)
		logic.ProcessError(err)
		Configurations = &data
	} else {
		if verbose {
			logger.Warn("JSON File is not installed!!")
		}
		Configurations = &entities.Configurations{
			SettingsName: "",
			Profiles:     []entities.Profile{},
		}
	}
}

func IsJsonFileInstalled() bool {
	return file.IsFile(JsonFile)
}

func GetVscodeListExtCommand(profileName string) models.Command {
	command := GetCodeCommand()
	command.Args = append(command.Args, profileName, "--list-extensions")
	command.Verbose = false
	return command
}

func GetCodeCommand() models.Command {
	command := models.Command{
		Cmd:      "code",
		Args:     []string{"--profile"},
		Verbose:  true,
		UseShell: false,
		IsThrow:  false,
	}
	return command
}

func OpenVscodeWithNewProfile(name string) {
	command := GetCodeCommand()
	command.Args = append(command.Args, name, "--wait")
	command.Verbose = true
	exe.ExecRealTime(command)
}
