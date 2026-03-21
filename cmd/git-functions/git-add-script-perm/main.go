package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-add-script-perm <script>",
		Short: "Add executable permission to a script",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(scriptArg string) {
	scriptFilename := file.Basename(scriptArg)
	updateCmd := models.Command{
		Cmd:     "git",
		Args:    []string{"update-index", "--chmod=+x", scriptArg},
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(updateCmd))
	lsCmd := models.Command{
		Cmd:     "git",
		Args:    []string{"ls-files", "--stage"},
		Verbose: true,
	}
	response, err := exe.Exec(lsCmd)
	logic.ProcessError(err)
	for line := range strings.SplitSeq(response, common.Eol()) {
		if strings.Contains(line, scriptFilename) {
			fmt.Println(line)
		}
	}
}

func main() {
	cobralib.Run()
}
