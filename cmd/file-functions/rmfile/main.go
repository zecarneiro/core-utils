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
		Use:   "rmfile [filepath]",
		Short: "Delete a file",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(filepath string) {
	if filepath == "." || filepath == ".." {
		filepath = ""
	}
	logic.ProcessError(file.DeleteFile(filepath))
}

func main() {
	cobralib.Run()
}
