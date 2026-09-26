package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"main/internal/dir"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var filter string

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "scoop-get-app-dir <app>",
		Short: "Get app dir by given app name. Ex: sublime-text",
		Args:  cobra.ExactArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(app string) {
	scoopDir := dir.GetScoopRootDir()
	appDir := file.JoinPath(scoopDir, "apps", app, "current")
	if !file.IsDir(scoopDir) {
		logic.ProcessError(fmt.Errorf("Not found scoop root dir!"))
	}
	if !file.IsDir(appDir) {
		logic.ProcessError(fmt.Errorf("Not found given app dir!"))
	}
	logger.Log(appDir)
}

func main() {
	cobralib.Run()
}
