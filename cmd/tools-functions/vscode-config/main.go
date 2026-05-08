package main

import (
	"golangutils/pkg/file"
	"golangutils/pkg/slice"

	"main/cmd/tools-functions/vscode-config/common"
	"main/internal/dir"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var processor *Processor

func init() {
	loadVars()
	setupCommand()
}

func loadVars() {
	common.JsonFile = file.JoinPath(dir.CoreUtilsUserConfig(), "vscode-config.json")
	processor = NewProcessor()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "vscode-config",
		Short: "VSCode Manager(Reset/Extract/Install)",
	}

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset VSCode",
		Run: func(cmd *cobra.Command, args []string) {
			processor.resetVscode()
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all profile by given JSON installed",
		Run: func(cmd *cobra.Command, args []string) {
			processor.listProfiles()
		},
	}

	profileExistCmd := &cobra.Command{
		Use:   "profile-exists <profile-name>",
		Short: "Check if given profile exists on VSCode",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileName := slice.ArrayToString(args)
			processor.profileStatus(profileName)
		},
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install JSON file and start install process(Optional)",
		Run: func(cmd *cobra.Command, args []string) {
			processor.install()
		},
	}
	installCmd.Flags().StringVarP(&processor.jsonFileArg, "json-file", "J", "", "Copy JSON File with all configs and profiles")
	installCmd.Flags().BoolVarP(&processor.settingsArg, "settings", "S", false, "Set settings from installed JSON file")
	installCmd.Flags().BoolVarP(&processor.processAllArg, "process-all", "P", false, "Process install flag for profiles, settings and extensions")

	extractCmd := &cobra.Command{
		Use:   "extract <profile-name>",
		Short: "Extract Data from JSON File",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			processor.extractProfileNameArg = slice.ArrayToString(args)
			processor.extract()
		},
	}
	extractCmd.Flags().BoolVarP(&processor.extractSettingsArg, "settings", "s", false, "Extract settings from installed JSON file")
	extractCmd.Flags().BoolVarP(&processor.extractExtsArg, "extensions", "e", false, "Extract all extensions by profile name from installed JSON file")
	extractCmd.Flags().BoolVarP(&processor.extractAllArg, "all", "a", false, "Extract all extensions and settings by profile name from installed JSON file")

	cobralib.CobraCmd.AddCommand(resetCmd, listCmd, profileExistCmd, installCmd, extractCmd)
}

func main() {
	cobralib.Run()
}
