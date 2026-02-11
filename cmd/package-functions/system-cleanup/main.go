package main

import (
	"fmt"
	"golangutils/pkg/logger"
	"golangutils/pkg/platform"
	"main/internal/libs"
)

func init() { setupCommand() }

func setupCommand() {}

func main() {
	fmt.Println("Not implemented yet!")
	libs.RunCoreUtilsCmd("error-log", "AAAA", "BBBB")
	//libs.RunCoreUtilsCmd("npm_clean")
	if platform.IsWindows() {
		libs.RunCoreUtilsCmd("scoop_clean")
	} else if platform.IsLinux() {
		logger.Info("AAAAAAA")
		/*for _, cmd := range []string{"apt_clean", "flatpak_clean", "deb_get_clean", "snap-clean"} {
		libs.RunCoreUtilsCmd(cmd)
		logic.Exit(0)
		}*/
	}
}
