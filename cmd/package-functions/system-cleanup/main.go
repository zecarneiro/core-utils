package main

import (
	"fmt"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"main/internal/generic"
)

func init() { setupCommand() }

func setupCommand() {}

func main() {
	fmt.Println("Not implemented yet!")
	generic.RunCoreUtilsCmd("npm_clean")
	if platform.IsWindows() {
		generic.RunCoreUtilsCmd("scoop_clean")
	} else if platform.IsLinux() {
		for _, cmd := range []string{"apt_clean", "flatpak_clean", "deb_get_clean", "snap-clean"} {
			generic.RunCoreUtilsCmd(cmd)
			logic.Exit(0)
		}
	}
}
