package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"slices"

	"main/cmd/tools-functions/winboat-volume/entities"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

const COMPOSE_WINDOWS_SERVICE_KEY = "windows"

var (
	removeArg bool

	dockerComposeDir string
	composeFile      string
	composeData      entities.Compose
)

func init() {
	console.EnableFeatures()
	if !system.IsAdmin() && platform.IsWindows() {
		logic.ProcessError(errors.New(system.NeedAdminAccessMsg))
	}
	loadData()
	setupCommand()
}

func loadData() {
	dockerComposeDir = fmt.Sprintf(`%s/.winboat`, system.HomeDir())
	composeFile = fmt.Sprintf(`%s/docker-compose.yml`, dockerComposeDir)
	if !file.IsFile(composeFile) {
		logic.ProcessError(fmt.Errorf(`Not found compose file: %s`, composeFile))
	}
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "winboat-volume <path>",
		Short: "Add path available on winboat",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().BoolVarP(&removeArg, "remove-path", "r", false, "Remove given path from docker compose")
	cobralib.WithRunArgsStr(process)
}

func getPathExposed(path string) string {
	return fmt.Sprintf(`/shared/%s`, file.Basename(path))
}

func existsPath(path string, services entities.Service) bool {
	return slices.Contains(services.Volumes, path)
}

func realoadCompose() {
	logger.Info("Reload Winboat VM")
	cmds := []string{"docker compose down", "docker compose up -d"}
	for _, cmd := range cmds {
		logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: cmd, Verbose: true, UseShell: true, Cwd: dockerComposeDir}))
	}
}

func loadComposeData() {
	logger.Info(fmt.Sprintf(`Reading Compose File: %s`, composeFile))
	data, err := file.ReadYamlFile[entities.Compose](composeFile)
	logic.ProcessError(err)
	composeData = data
}

func updateComposeFile(services entities.Service) {
	composeData.Services[COMPOSE_WINDOWS_SERVICE_KEY] = services
	logic.ProcessError(file.WriteYamlFile(composeFile, composeData))
}

func process(path string) {
	canReloadCompose := false
	if !str.IsEmpty(path) && file.IsDir(path) {
		loadComposeData()
		services := composeData.Services[COMPOSE_WINDOWS_SERVICE_KEY]
		exposedPath := getPathExposed(path)
		volumePath := fmt.Sprintf(`%s:%s`, path, exposedPath)
		if !existsPath(volumePath, services) {
			if !removeArg {
				services.Volumes = append(services.Volumes, volumePath)
				logger.Info(fmt.Sprintf(`Exposing given path on: %s`, exposedPath))
				updateComposeFile(services)
				canReloadCompose = true
			}
		} else {
			if removeArg {
				services.Volumes = slice.FilterArray(services.Volumes, func(volume string) bool {
					return volume != volumePath
				})
				logger.Info(fmt.Sprintf(`Remove Exposed given path: %s`, path))
				updateComposeFile(services)
				canReloadCompose = true
			} else {
				logger.Ok("Already exist given path.")
			}
		}
	} else {
		logger.ErrorStr("Invalid given path.")
	}
	if canReloadCompose {
		realoadCompose()
		logger.Ok("Done.")
	}
}

func main() {
	cobralib.Run()
}
