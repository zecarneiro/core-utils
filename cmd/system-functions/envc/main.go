package main

import (
	"fmt"
	"golangutils/pkg/env"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() {
	load()
	setupCommand()
}

func load() {
	envManager = libs.NewEnvManager()
	envManager.SetSystemConfig()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   fmt.Sprintf("%s <command>", appName),
		Short: "Manage ENV variables",
	}
	if !platform.IsWindows() {
		cobralib.CobraCmd.Long = fmt.Sprintf(warnningHelp, envManager.EnvFile)
	}

	// Exists
	existsCmd := &cobra.Command{
		Use:   "exists <name>",
		Short: "Check if env var exists",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			exists(slice.ArrayToString(args))
		},
	}

	// Exists Value
	existsValueCmd := &cobra.Command{
		Use:   "exists-value",
		Short: "Check if value of env var exists",
		Run: func(cmd *cobra.Command, args []string) {
			name, errName := cmd.Flags().GetString("name")
			value, errValue := cmd.Flags().GetString("value")
			logic.ProcessError(errName)
			logic.ProcessError(errValue)
			existsValue(name, value)
		},
	}
	existsValueCmd.Flags().StringP("name", "n", "", "Name of env var")
	logic.ProcessError(existsValueCmd.MarkFlagRequired("name"))
	existsValueCmd.Flags().StringP("value", "v", "", "Value of given name var env")
	logic.ProcessError(existsValueCmd.MarkFlagRequired("value"))

	// Add
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add env var with value",
		Long:  "Add env var with value. Work onnly for User",
		Run: func(cmd *cobra.Command, args []string) {
			name, errName := cmd.Flags().GetString("name")
			value, errValue := cmd.Flags().GetString("value")
			isReplaceAll, errReplaceAll := cmd.Flags().GetBool("replace-all")
			logic.ProcessError(errName)
			logic.ProcessError(errValue)
			logic.ProcessError(errReplaceAll)
			values := env.ConvValuesArr(value)
			validateValue(values, true)
			processAddArr(name, values, isReplaceAll)
			logger.Ok(fmt.Sprintf("Added env %s with value %s", name, value))
			warnningMsg()
		},
	}
	addCmd.Flags().StringP("name", "n", "", "Name of env var")
	logic.ProcessError(addCmd.MarkFlagRequired("name"))
	addCmd.Flags().StringP("value", "v", "", "Value of given name var env")
	logic.ProcessError(addCmd.MarkFlagRequired("value"))
	addCmd.Flags().BoolP("replace-all", "r", false, "Replace all values of given env for a new given value.")

	// Delete duplicated
	cleanCmd := &cobra.Command{
		Use:   "clean <name>",
		Short: "Delete duplicated value of env var. Work onnly for User",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			processClean(slice.ArrayToString(args))
			warnningMsg()
		},
	}

	// Delete
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete env var or delete env value of given env var. Work onnly for User",
		Run: func(cmd *cobra.Command, args []string) {
			name, errName := cmd.Flags().GetString("name")
			value, errValue := cmd.Flags().GetString("value")
			logic.ProcessError(errName)
			logic.ProcessError(errValue)
			processDelete(name, value)
			logger.Ok(fmt.Sprintf("Deleted env %s with value %s", name, value))
			warnningMsg()
		},
	}
	deleteCmd.Flags().StringP("name", "n", "", "Name of env var")
	logic.ProcessError(deleteCmd.MarkFlagRequired("name"))
	deleteCmd.Flags().StringP("value", "v", "", "Value of given name var env")

	// Set
	cobralib.CobraCmd.AddCommand(existsCmd, existsValueCmd, addCmd, cleanCmd, deleteCmd)
}

func warnningMsg() {
	logger.Info("Done. Please, run this command: reload-shell")
	logger.Warn("Check if take efect. If not, restart all terminal.")
}

func main() {
	cobralib.Run()
}
