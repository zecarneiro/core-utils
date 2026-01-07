package main

import (
	"golangutils/pkg/file"
	"golangutils/pkg/logic"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "mkdirc [directory_path]",
		Short: "Create a directory if it doesn't exist",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(path string) {
	if path == "." || path == ".." {
		path = ""
	}
	logic.ProcessError(file.CreateDirectory(path, true))
}

func main() {
	cobralib.Run()
}
