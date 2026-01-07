package main

import (
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var force bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "rmdirc [directory_path]",
		Short: "Delete a directory",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().BoolVarP(&force, "force", "f", false, "Force delete operation")
	cobralib.WithRunArgsStr(process)
}

func process(path string) {
	path = logic.Ternary(path == "." || !libs.IsValidPathArg(path), "", path)
	filesInfo, err := file.ReadDirRecursive(path)
	logic.ProcessError(err)
	if !force && (len(filesInfo.Directories) > 0 || len(filesInfo.Files) > 0) && !console.Confirm("Directory is not empty. Conitinue?", true) {
		logic.Exit(0)
	}
	logic.ProcessError(file.DeleteDirectory(path))
}

func main() {
	cobralib.Run()
}
