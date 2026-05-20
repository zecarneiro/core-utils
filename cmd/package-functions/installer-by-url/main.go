package main

import (
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/slice"

	"main/cmd/package-functions/installer-by-url/args"
	"main/internal/dir"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() {
	loadData()
	setupCommand()
}

func loadData() {
	configFileDir = file.JoinPath(dir.CoreUtilsUserConfig(), "install-by-url")
	if !file.IsDir(configFileDir) {
		logic.ProcessError(file.CreateDirectory(configFileDir, true))
	}
	destScriptDir = file.JoinPath(dir.CoreUtilsLocal(), "installer-by-url")
	if !file.IsDir(destScriptDir) {
		logic.ProcessError(file.CreateDirectory(destScriptDir, true))
	}
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "installer-by-url",
		Short: "Install files from URL",
		Args:  cobra.NoArgs,
	}

	// Github
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Use github to use to get the app version and upgrade the app.",
		Run: func(cmd *cobra.Command, args []string) {
			github(cmd)
		},
	}
	args.BuildGithubFlags(githubCmd)

	// Generic
	genericCmd := &cobra.Command{
		Use:   "generic",
		Short: "Only download the file from URL",
		Run: func(cmd *cobra.Command, args []string) {
			generic(cmd)
		},
	}
	args.BuildGenericFlags(genericCmd)

	// List
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all installers",
		Run: func(cmd *cobra.Command, args []string) {
			listAll()
		},
	}

	// Run
	runCmd := &cobra.Command{
		Use:   "run [APP_NAME| ]",
		Short: "Run installer. If name not set, will run all apps",
		Run: func(cmd *cobra.Command, args []string) {
			run(slice.ArrayToString(args))
		},
	}
	runCmd.Flags().BoolVarP(&forceInstallArg, "force", "f", false, "Force install")

	// Uninstall
	uninstallCmd := &cobra.Command{
		Use:   "uninstall [APP_NAME| ]",
		Short: "Uninstall installer",
		Run: func(cmd *cobra.Command, args []string) {
			uninstall(slice.ArrayToString(args))
		},
	}

	// Set
	cobralib.CobraCmd.AddCommand(githubCmd, genericCmd, listCmd, runCmd, uninstallCmd)
}

func main() {
	cobralib.Run()
}
