package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "file-encoding [filepath]",
		Short: "Get encoding of file",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(filepath string) {
	if filepath == "." || filepath == ".." {
		filepath = ""
	}
	if !file.IsFile(filepath) {
		logic.ProcessError(fmt.Errorf("%s is not a file or file not found", filepath))
	}
	name, err := file.GetFileEncoding(filepath)
	logic.ProcessError(err)
	fmt.Println(name)
}

func main() {
	cobralib.Run()
}
