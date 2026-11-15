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
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "create-boot-app-entry",
		Short: "Create a Boot APP for system",
	}
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "VBS filename (required)")
	cobralib.CobraCmd.Flags().StringVarP(&commandArg, "command", "c", "", "Command to run or binary (required)")
	cobralib.WithRun(process)
}

func process() {
	outputPath, err := system.CreateBootEntry(nameArg, commandArg)
	logic.ProcessError(err)
	fmt.Printf("Created: %s%s", outputPath, common.Eol())
}

func main() {
	cobralib.Run()
}
