package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"regexp"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	nameArg    string
	commandArg string
	waitArg    bool

	template = `Dim WinScriptHost
Set WinScriptHost = CreateObject("WScript.Shell")
WinScriptHost.Run """%s""", %d, False
Set WinScriptHost = Nothing
`
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

func validateName(name string) {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if str.IsEmpty(name) || !re.MatchString(name) {
		logic.ProcessError(fmt.Errorf("Invalid given name!. Name must not be empty and accept only A-Z, a-z, 0-9, - and _"))
	}
}

func process() {
	validateName(nameArg)
	if str.IsEmpty(commandArg) {
		logic.ProcessError(fmt.Errorf("Command is required"))
	}
	waitInt := 0
	if waitArg {
		waitInt = 1
	}
	outputPath := file.JoinPath(cobralib.GetWorkingDir(), fmt.Sprintf("%s.vbs", nameArg))
	content := fmt.Sprintf(template, commandArg, waitInt)
	fileConfig := models.FileWriterConfig{
		File:        outputPath,
		Data:        content,
		IsAppend:    false,
		WithUtf8BOM: false,
	}
	logic.ProcessError(file.WriteFile(fileConfig))
	fmt.Printf("Created: %s%s", outputPath, common.Eol())
}

func main() {
	cobralib.Run()
}
