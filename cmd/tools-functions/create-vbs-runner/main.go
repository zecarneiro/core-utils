package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/logic"
	"golangutils/pkg/system"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	nameArg    string
	commandArg string
	waitArg    bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "create-vbs-runner",
		Short: "Create a VBS script to run a command hidden",
		Long:  "Creates a VBS file that runs a command without showing a window",
	}
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "VBS filename (required)")
	cobralib.CobraCmd.Flags().StringVarP(&commandArg, "command", "c", "", "Command to run (required)")
	cobralib.CobraCmd.Flags().BoolVarP(&waitArg, "keep-window", "k", false, "Wait for command to finish")
	cobralib.WithWorkingDirDefault()
	cobralib.WithRun(process)
}

func process() {
	outputPath, err := system.CreateVisualBasicScript(nameArg, commandArg, waitArg, cobralib.GetWorkingDir())
	logic.ProcessError(err)
	fmt.Printf("Created: %s%s", outputPath, common.Eol())
}

func main() {
	cobralib.Run()
}
