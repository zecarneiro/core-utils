package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"slices"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	desktopEnv         enums.DesktopEnvType
	fullPath           string
	linuxBookmarksFile string
	linuxBookmarksData string
	dataGtkTemplate    = "file://%s"
	dataKdeTemplate    = `
<bookmark href="file://%s">
  <title>%s</title>
</bookmark>
`
	cmdWindowsCheckTemplate = "$shell = New-Object -ComObject Shell.Application; $shell.Namespace(\"shell:::{679f85cb-0220-4080-b29b-5540cc05aab6}\").Items() | Where-Object { $_.Path -eq \"%s\" }"
	cmdWindowsSetTemplate   = "$shell = New-Object -ComObject Shell.Application; $folder = $shell.Namespace(\"%s\"); $folder.Self.InvokeVerb(\"pintohome\")"
)

func init() {
	desktopEnv = system.GetDesketopEnv()
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "pin-into-file-explorer-bookmarks <directory>",
		Short: "Pin File/Dir into your file explorer(Only working and tested for Microsoft file explorer, Nautilus)",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func validatePlatform() {
	validDesktopEnv := []enums.DesktopEnvType{enums.GnomeDE, enums.KdeDE, enums.WindowsDE}
	if !slices.Contains(validDesktopEnv, desktopEnv) {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
}

func getFullPath(pathArg string) string {
	fullPath, err := file.GetFullPath(pathArg)
	logic.ProcessError(err)
	return fullPath
}

func getLinuxBookmarksFile() string {
	if platform.IsLinux() {
		if desktopEnv.Equals(enums.GnomeDE) {
			return file.JoinPath(system.HomeDir(), ".config/gtk-3.0/bookmarks")
		} else if desktopEnv.Equals(enums.KdeDE) {
			return file.JoinPath(system.HomeDir(), ".config/gtk-3.0/bookmarks")
		}
	}
	return ""
}

func getLinuxBookmarksData(fullPath string) string {
	if desktopEnv.Equals(enums.GnomeDE) {
		return fmt.Sprintf(dataGtkTemplate, fullPath)
	} else if desktopEnv.Equals(enums.KdeDE) {
		return fmt.Sprintf(dataKdeTemplate, fullPath, file.FileName(fullPath))
	}
	return ""
}

func alreadyExist() bool {
	if platform.IsWindows() {
		res, err := exe.Exec(models.Command{Cmd: fmt.Sprintf(cmdWindowsCheckTemplate, fullPath), UseShell: true})
		if err != nil {
			return false
		}
		return logic.Ternary(str.IsEmpty(res), false, true)
	} else if desktopEnv.Equals(enums.KdeDE) {
		alreadySet, err := file.FileTextContains(linuxBookmarksFile, linuxBookmarksData, false)
		if err != nil {
			alreadySet = false
		}
		return alreadySet
	} else if desktopEnv.Equals(enums.GnomeDE) {
		alreadySet := false
		file.ReadFileLineByLine(linuxBookmarksFile, func(line string) {
			if !alreadySet && linuxBookmarksData == line {
				alreadySet = true
			}
		})
		return alreadySet
	}
	return false
}

func process(args string) {
	validatePlatform()
	if str.IsEmpty(args) || !file.IsDir(args) {
		logic.ProcessError(fmt.Errorf("Invalid given file: %s", args))
	}
	fullPath = getFullPath(args)
	linuxBookmarksFile = getLinuxBookmarksFile()
	linuxBookmarksData = getLinuxBookmarksData(fullPath)
	if alreadyExist() {
		logger.Info(fmt.Sprintf("Path already set as bookmarks/quick access: %s", fullPath))
	} else {
		logger.Info(fmt.Sprintf("Set path on bookmarks/quick access: %s", fullPath))
		if platform.IsWindows() {
			logger.Error(exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf(cmdWindowsSetTemplate, fullPath), UseShell: true}))
		} else {
			fileWriter := models.FileWriterConfig{
				File:        linuxBookmarksFile,
				Data:        linuxBookmarksData,
				IsAppend:    true,
				IsCreateDir: true,
				WithUtf8BOM: false,
			}
			logger.Error(file.WriteFile(fileWriter))
		}
	}
}

func main() {
	cobralib.Run()
}
