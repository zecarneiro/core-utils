package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	filterArg          string
	caseInsensitiveArg bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "start-apps",
		Short: "List installed applications",
	}
	cobralib.CobraCmd.Flags().StringVarP(&filterArg, "filter", "f", "", "Filter by application name")
	cobralib.CobraCmd.Flags().BoolVarP(&caseInsensitiveArg, "case-insensitive", "i", false, "Enable filter with case insensitive")
	cobralib.WithRun(process)
}

func process() {
	if platform.IsWindows() {
		processWindows()
	} else {
		processLinux()
	}
}

func processWindows() {
	cmd := models.Command{
		Cmd:        "Get-StartApps",
		UseShell:   true,
		Verbose:    true,
		ShellToUse: enums.PowerShell,
	}
	response, err := exe.Exec(cmd)
	logic.ProcessError(err)
	for line := range strings.SplitSeq(response, common.Eol()) {
		if !str.IsEmpty(filterArg) && str.Contains(line, filterArg, caseInsensitiveArg) {
			fmt.Println(line)
		} else {
			fmt.Println(line)
		}
	}
}

func processLinux() {
	desktopDirs := []string{
		"/usr/share/applications",
		"/var/lib/snapd/desktop/applications",
		"/var/lib/flatpak/exports/share/applications",
		file.JoinPath(system.HomeDir(), ".local/share/applications"),
		file.JoinPath(system.HomeDir(), ".local/share/flatpak/exports/share/applications"),
	}
	for _, dir := range desktopDirs {
		if !file.IsDir(dir) {
			continue
		}
		filesList, err := file.ReadDirWithFilter(dir, "*.desktop")
		if err != nil {
			logger.Error(err)
			continue
		}
		for _, filePath := range filesList {
			var name, exec, icon string
			content, err := file.ReadFile(filePath)
			logger.Error(err)
			lines := strings.Split(content, common.Eol())
			for _, line := range lines {
				if strings.HasPrefix(line, "Name=") {
					name = strings.TrimPrefix(line, "Name=")
				} else if strings.HasPrefix(line, "Exec=") {
					exec = strings.TrimPrefix(line, "Exec=")
				} else if strings.HasPrefix(line, "Icon=") {
					icon = strings.TrimPrefix(line, "Icon=")
				}
			}
			if !str.IsEmpty(name) && !str.IsEmpty(exec) {
				canPrint := true
				if !str.IsEmpty(filterArg) {
					if !str.Contains(name, filterArg, caseInsensitiveArg) {
						canPrint = false
					}
				}
				if canPrint {
					fmt.Println(name)
					logger.Log(fmt.Sprintf("  Exec: %s", exec))
					logger.Log(fmt.Sprintf("  Icon: %s", icon))
					logger.Log(fmt.Sprintf("  File: %s", filePath))
					fmt.Println()
				}
			}
		}
	}
}

func main() {
	cobralib.Run()
}
