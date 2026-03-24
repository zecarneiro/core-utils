package main

import (
	"golangutils/pkg/platform"
	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "system-cleanup",
		Short: "Clean system by removing all non necessary packages",
	}
	cobralib.WithRun(process)
}

func process() {
	if platform.IsWindows() {
		libs.RunCoreUtilsCmd("scoop-clean", true)
	} else if platform.IsLinux() {
		for _, cmd := range []string{"apt-clean", "flatpak-clean", "deb-get-clean", "snap-clean"} {
			libs.RunCoreUtilsCmd(cmd, true)
		}
	}
}

func main() {
	cobralib.Run()
}
