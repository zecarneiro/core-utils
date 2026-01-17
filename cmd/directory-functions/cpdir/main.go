package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"golangutils/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	source      string
	destination string
	force       bool
	rootCmd     *cobra.Command
)

func init() {
	sourceKey := "source"
	destinationKey := "destination"
	rootCmd = &cobra.Command{
		Use:   "cpdir",
		Short: "Copy directory to another location",
	}
	rootCmd.Flags().StringVarP(&source, sourceKey, "s", "", "Source directory")
	rootCmd.MarkFlagRequired(sourceKey)
	rootCmd.Flags().StringVarP(&destination, destinationKey, "d", "", "Destination directory")
	rootCmd.MarkFlagRequired(destinationKey)
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "Force operation, even if destination already exists")
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		currentDir, err := file.GetCurrentDir()
		generic.ProcessError(err)
		if source == "" || destination == "" || source == ".." || destination == ".." {
			logger.Error("Source and/or destination must be provided")
			os.Exit(1)
		} else if source == destination {
			logger.Error("Source and destination cannot be the same")
			os.Exit(1)
		}
		if !force && file.IsDir(destination) && !generic.Confirm("Destination already exist. Conitinue?", true) {
			os.Exit(0)
		}
		source = common.Ternary(source == ".", currentDir, source)
		destination = common.Ternary(destination == ".", currentDir, destination)
		generic.ProcessError(file.CopyDir(source, destination))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
