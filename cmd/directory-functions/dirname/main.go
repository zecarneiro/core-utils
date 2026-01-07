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
		Use:   "dirname [path]",
		Short: "Get dirname of given path",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(path string) {
	if path == "." {
		path = libs.GetCurrentDir(true)
	}
	fmt.Println(file.Dirname(path))
}

func main() {
	cobralib.Run()
}
