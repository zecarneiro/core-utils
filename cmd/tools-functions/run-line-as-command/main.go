package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"strings"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	fileArg       string
	sectionArg    string
	allArg        bool
	dataBySection []string

	sectionAll = "ALL"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "read-section",
		Short: "Read a specific section from a file",
		Long:  "Reads content between [SECTION] tags",
	}
	cobralib.CobraCmd.Flags().StringVarP(&fileArg, "file", "f", "", "File to read (required)")
	cobralib.CobraCmd.Flags().StringVarP(&sectionArg, "section", "s", "", "Section name (e.g., WIN, LINUX, ALL)")
	cobralib.CobraCmd.Flags().BoolVarP(&allArg, "all", "a", false, "Include ALL section with the selected section")
	cobralib.WithRun(process)
}

func process() {
	fileArg = file.ResolvePath(fileArg)
	if str.IsEmpty(fileArg) || !file.IsFile(fileArg) {
		logic.ProcessError(errors.New("Please give a valid file argument"))
	}
	if str.IsEmpty(sectionArg) {
		logic.ProcessError(errors.New("Please give a valid section argument"))
	}
	sectionArg = strings.ToUpper(sectionArg)
	shouldInsert := false
	err := file.ReadFileLineByLine(fileArg, func(line string) {
		isSectionLine := false
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			isSectionLine = true
			sectionName := strings.ToUpper(strings.Trim(line, "[]"))
			sectionName = strings.ToUpper(sectionName)
			if sectionName == sectionArg || (allArg && sectionName == sectionAll) {
				shouldInsert = true
			} else {
				shouldInsert = false
			}
		}
		if shouldInsert && !isSectionLine {
			dataBySection = append(dataBySection, line)
		}
	})
	logic.ProcessError(err)
	dataBySectionStr := slice.ObjArrayToStringBySep(dataBySection, common.Eol())
	tempFile := system.GenerateTempFile(file.FileName(fileArg))
	if platform.IsWindows() {
		tempFile = fmt.Sprintf("%s.ps1", tempFile)
	}
	fileConfig := models.FileWriterConfig{
		File:        tempFile,
		Data:        dataBySectionStr,
		IsAppend:    false,
		WithUtf8BOM: false,
		IsCreateDir: true,
	}
	logic.ProcessError(file.WriteFile(fileConfig))
	shellType := logic.Ternary(platform.IsWindows(), enums.PowerShell, shell.GetCurrentShellSimple())
	libs.RunCoreUtilsCmd("run-shell-script", false, "-f", tempFile, "-s", shellType.String())
	logic.ProcessError(file.DeleteFile(tempFile))
}

func main() {
	cobralib.Run()
}
