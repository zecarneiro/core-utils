package main

import (
	"golangutils/pkg/slice"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() {
	loadAndValidateVars()
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   appName,
		Short: "Manage script updaters",
	}

	// Install
	installCmd := &cobra.Command{
		Use:   "install <filepath>",
		Short: "Install shell script by given file script",
		Args:  cobra.MinimumNArgs(1),
		Run:   install,
	}
	installCmd.Flags().BoolP("is-package-manager", "p", false, "Set if script is package manager")
	installCmd.Flags().BoolP("force", "f", false, "Force install")

	// Uninstall
	uninstallCmd := &cobra.Command{
		Use:   "uninstall <script_name>",
		Short: "Uninstall shell script by given file script",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			uninstall(slice.ArrayToString(args))
		},
	}

	// Run
	runCmd := &cobra.Command{
		Use:   "run [ |script_name]",
		Short: "Process specific script package manager",
		Run: func(cmd *cobra.Command, args []string) {
			run(slice.ArrayToString(args))
		},
	}

	// List
	listCmd := &cobra.Command{
		Use:   "list [ |filter|]",
		Short: "List installed scripts",
		Run: func(cmd *cobra.Command, args []string) {
			list(slice.ArrayToString(args))
		},
	}

	// Update
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update all scripts for others shell's",
		Run: func(cmd *cobra.Command, args []string) {
			update()
		},
	}

	// Script Version
	scriptVersionCmd := &cobra.Command{
		Use:   "script-version <script>",
		Short: "Version of specific script package manager",
		Run:   scriptVersionManager,
	}
	scriptVersionCmd.Flags().StringP("version", "v", "", "Script version")

	// Set
	cobralib.CobraCmd.AddCommand(installCmd, uninstallCmd, runCmd, listCmd, updateCmd, scriptVersionCmd)
}

func main() {
	cobralib.Run()
}
