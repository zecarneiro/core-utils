package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"

	"main/internal/libs"
	"main/internal/libs/cobralib"
	"main/internal/libs/golangutilslib"

	"github.com/spf13/cobra"
)

var (
	source      string
	destination string
	force       bool
)

func init() { setupCommand() }

func setupCommand() {
	sourceKey := "source"
	destinationKey := "destination"
	cobralib.CobraCmd = &cobra.Command{
		Use:   "cpdir",
		Short: "Copy directory to another location",
	}
	cobralib.CobraCmd.Flags().StringVarP(&source, sourceKey, "s", "", "Source directory")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired(sourceKey))
	cobralib.CobraCmd.Flags().StringVarP(&destination, destinationKey, "d", "", "Destination directory")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired(destinationKey))
	cobralib.CobraCmd.Flags().BoolVarP(&force, "force", "f", false, "Force operation, even if destination already exists")
	cobralib.WithRun(process)
}

func process() {
	currentDir := libs.GetCurrentDir(true)
	if source == "" || destination == "" || source == ".." || destination == ".." {
		logic.ProcessError(fmt.Errorf("Source and/or destination must be provided"))
	} else if source == destination {
		logic.ProcessError(fmt.Errorf("Source and destination cannot be the same"))
	}
	if !force && file.IsDir(destination) && !console.Confirm("Destination already exist. Conitinue?", true) {
		logic.Exit(0)
	}
	source = logic.Ternary(source == ".", currentDir, source)
	destination = logic.Ternary(destination == ".", currentDir, destination)
	logic.ProcessError(golangutilslib.FuncCopyDir(source, destination))
}

func main() {
	cobralib.Run()
}
