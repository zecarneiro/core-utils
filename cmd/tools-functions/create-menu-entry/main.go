package main

import (
	"fmt"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "create-menu-entry",
		Short: "Create a menu entry (.desktop file on Linux)",
		Long:  "Creates a menu entry in the system",
	}
	cobralib.CobraCmd.Flags().StringVarP(&menuEntryData.name, "name", "n", "", "Application name (required)")
	cobralib.CobraCmd.Flags().StringVarP(&menuEntryData.exec, "exec", "e", "", "Executable path (required)")
	cobralib.CobraCmd.Flags().StringArrayVarP(&menuEntryData.args, "args", "a", []string{}, "Arguments to pass to the executable")
	cobralib.CobraCmd.Flags().StringVarP(&menuEntryData.icon, "icon", "i", "", "Icon path")
	cobralib.CobraCmd.Flags().StringVarP(&menuEntryData.categories, "categories", "c", "", "Categories (e.g., Utility;Development)")
	cobralib.CobraCmd.Flags().StringVarP(&menuEntryData.comment, "comment", "m", "", "Comment/description")
	cobralib.CobraCmd.Flags().BoolVarP(&menuEntryData.terminal, "terminal", "t", false, "Run in terminal")
	if platform.IsWindows() {
		cobralib.CobraCmd.Flags().BoolVarP(&menuEntryData.runAsAdmin, "admin", "A", false, "Run as administrator (Windows SO only)")
	}
	cobralib.WithRun(process)
}

func process() {
	if str.IsEmpty(menuEntryData.name) {
		logic.ProcessError(fmt.Errorf("Name is required"))
	}
	if str.IsEmpty(menuEntryData.exec) {
		logic.ProcessError(fmt.Errorf("Exec is required"))
	}
	menuEntryData.argsStr = slice.ArrayToString(slice.RemoveAllEmpty(menuEntryData.args))
	start()
}

func main() {
	cobralib.Run()
}
