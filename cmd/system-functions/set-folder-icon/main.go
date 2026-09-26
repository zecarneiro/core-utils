package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	folder string
	icon   string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "set-folder-icon",
		Short: "Set folder icon",
	}
	cobralib.CobraCmd.Flags().StringVarP(&folder, "folder", "f", "", "Full folder path (required)")
	cobralib.CobraCmd.Flags().StringVarP(&icon, "icon", "i", "", "Full icon path (required)")
	cobralib.WithRun(process)
}

func windowsProcess() {
	desktopIniName := "desktop.ini"
	desktopIniFile := file.JoinPath(folder, desktopIniName)
	desktopIniContent := `
	[.ShellClassInfo]
	IconResource=%s,0
	`
	logic.ProcessError(file.WriteFile(models.FileWriterConfig{File: desktopIniFile, Data: fmt.Sprintf(desktopIniContent, icon), IsAppend: false, IsCreateDir: false}))
	libs.RunCoreUtilsCmd("restore-folder-icon", false, folder)
}

func getFullPath(relativePath string) string {
	fullPath, err := file.GetFullPath(relativePath)
	logic.ProcessError(err)
	return fullPath
}

func process() {
	isInvalidPlatform := false
	folder = file.ResolvePath(folder)
	icon = file.ResolvePath(icon)
	if str.IsEmpty(folder) || !file.IsDir(folder) {
		logic.ProcessError(fmt.Errorf("Invalid given folder!"))
	}
	if str.IsEmpty(icon) || !file.IsFile(icon) {
		logic.ProcessError(fmt.Errorf("Invalid given icon!"))
	}
	folder = getFullPath(folder)
	icon = getFullPath(icon)
	if platform.IsWindows() {
		windowsProcess()
	} else if platform.IsLinux() {
		if console.CmdExists("gio") {
			logger.Error(exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf(`gio set "%s" metadata::custom-icon "file://%s"`, folder, icon), Verbose: true}))
		} else {
			isInvalidPlatform = true
		}
	} else {
		isInvalidPlatform = true
	}
	if isInvalidPlatform {
		logic.ProcessError(fmt.Errorf("%s", platform.UnsupportedMSG))
	}
}

func main() {
	cobralib.Run()
}
