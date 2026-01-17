package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	force   bool
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "rmdirc [directory_path]",
		Short: "Delete a directory",
		Args:  cobra.MinimumNArgs(1),
	}
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "Force delete operation")
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		data := strings.Join(args, " ")
		data = common.Ternary(data == "." || data == "..", "", data)
		filesInfo, err := file.ReadDirRecursive(data)
		generic.ProcessError(err)
		if !force && (len(filesInfo.Directories) > 0 || len(filesInfo.Files) > 0) && !generic.Confirm("Directory is not empty. Conitinue?", true) {
			os.Exit(0)
		}
		generic.ProcessError(file.DeleteDirectory(data))
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
