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
		Use:   "direxists [directory_path]",
		Short: "Check if a directory exists",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(path string) {
	if path == "." || path == ".." {
		path = ""
	}
	fmt.Println(file.IsDir(path))
}

func main() {
	cobralib.Run()
}
