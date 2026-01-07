package main

import (
	"fmt"
	"golangutils/pkg/file"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "fileexists [filepath]",
		Short: "Check if a file exists",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(filepath string) {
	if filepath == "." || filepath == ".." {
		filepath = ""
	}
	fmt.Println(file.IsFile(filepath))
}

func main() {
	cobralib.Run()
}
