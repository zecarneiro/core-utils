package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-repo-backup <url>",
		Short: "Clone a git repository as mirror (backup)",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(url string) {
	if str.IsEmpty(url) {
		logic.ProcessError(fmt.Errorf("invalid given URL"))
	}
	cmd := models.Command{
		Cmd:     "git",
		Args:    []string{"clone", "--mirror", url},
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
