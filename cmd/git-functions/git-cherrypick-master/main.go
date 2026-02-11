package main

import (
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "git-cherrypick-master <commit>...",
		Short: "Cherry-pick commits to master branch",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(args string) {
	cmd := models.Command{
		Cmd:     "git",
		Args:    append([]string{"cherry-pick", "-m", "1"}, args),
		Verbose: true,
	}
	logic.ProcessError(exe.ExecRealTime(cmd))
}

func main() {
	cobralib.Run()
}
