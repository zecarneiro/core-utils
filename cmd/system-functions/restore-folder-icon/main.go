package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "restore-folder-icon <folder-full-path>",
		Short: "Restore folder icon. On windows the folder must contain desktop.ini file",
		Args:  cobra.ExactArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func getFullFolderPath(folder string) string {
	fullPath, err := file.GetFullPath(folder)
	logic.ProcessError(err)
	return fullPath
}

func process(folder string) {
	folder = file.ResolvePath(folder)
	if str.IsEmpty(folder) || !file.IsDir(folder) {
		logic.ProcessError(fmt.Errorf("Invalid given folder!"))
	}
	folder = getFullFolderPath(folder)
	if platform.IsWindows() {
		desktopIniName := "desktop.ini"
		if !file.IsFile(file.JoinPath(folder, desktopIniName)) {
			logic.ProcessError(fmt.Errorf("Not found %s file!", desktopIniName))
		}
		logger.Error(exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf("attrib.exe +s \"%s\"", file.Basename(folder)), Verbose: true, Cwd: file.Dirname(folder)}))
		logger.Error(exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf("attrib.exe +h +s %s", desktopIniName), Verbose: true, Cwd: folder}))
	} else {
		logic.ProcessError(fmt.Errorf("%s", platform.UnsupportedMSG))
	}
}

func main() {
	cobralib.Run()
}
