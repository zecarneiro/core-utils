package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	filter     string
	includeAll bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "appimage-list",
		Short: "Get list of all packages or by filter",
	}
	cobralib.CobraCmd.Flags().StringVarP(&filter, "filter", "f", "", "Package to search")
	cobralib.CobraCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "Include all informations")
	cobralib.WithRun(process)
}

func printAllData(data libs.GearlevelInstalledAppInfo, registerName string) {
	logger.Separator()
	fmt.Printf("Name: %s%s", data.Name, common.Eol())
	fmt.Printf("Register Name: %s%s", registerName, common.Eol())
	fmt.Printf("Appimage File: %s%s", data.AppimageFile, common.Eol())
	fmt.Printf("Version: %s%s", data.Version, common.Eol())
}

func printData(data libs.GearlevelInstalledAppInfo, registerName string) {
	if includeAll {
		printAllData(data, registerName)
	} else {
		fmt.Println(registerName)
	}
}

func process() {
	hasFilter := logic.Ternary(!str.IsEmpty(filter), true, false)
	installedApps := libs.GetGearlevelListOfInstalledApps().Installed
	for _, app := range installedApps {
		registerName := app.GetRegisterName()
		if hasFilter {
			if str.Contains(registerName, filter, true) {
				printData(app, registerName)
			}
		} else {
			printData(app, registerName)
		}
	}
}

func main() {
	cobralib.Run()
}
