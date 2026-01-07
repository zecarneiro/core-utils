package main

import (
	"fmt"
	"golangutils/pkg/file"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "basename [filepath]",
		Short: "Get basename of given path",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(filepath string) {
	if filepath == "." {
		filepath = libs.GetCurrentDir(true)
	} else if !libs.IsValidPathArg(filepath) {
		filepath = ""
	}
	fmt.Println(file.Basename(filepath))
}

func main() {
	cobralib.Run()
}
