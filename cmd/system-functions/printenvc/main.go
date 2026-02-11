package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/env"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/str"
	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	list         bool
	listNamesArg bool
	variableArg  string
	envManager   *libs.EnvManager
)

func init() {
	envManager = libs.NewEnvManager()
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "printenvc",
		Short: "List all env vars names with their values",
	}
	cobralib.CobraCmd.Flags().BoolP("list", "l", false, "List all env's vars with theirs values (Selected by default)")
	cobralib.CobraCmd.Flags().BoolVarP(&listNamesArg, "list-names", "n", false, "List all names of env's vars")
	cobralib.CobraCmd.Flags().StringVarP(&variableArg, "list-values", "v", "", "List all values of given env vars")
	cobralib.WithRun(process)
}

func warnningMsg() {
	logger.Warn("If you get some wrong values or names, you may need to reload your shell!")
}

func printValues(name string) {
	if str.IsEmpty(name) {
		logic.ProcessError(fmt.Errorf("Invalid given name!"))
	}
	values := env.Get(name)
	fmt.Printf("%s has %d entries:%s", name, len(values), common.Eol())
	for index, value := range values {
		fmt.Printf("%d. %s%s", index+1, value, common.Eol())
	}
}

func listNames() {
	warnningMsg()
	values := env.ListFullInfo()
	count := 1
	for name := range values {
		output := fmt.Sprintf("%d. %s", count, name)
		fmt.Println(output)
		count++
	}
}

func listFullInfo() {
	warnningMsg()
	values := env.ListFullInfo()
	logger.WithSeparatorLength(20)
	for name := range values {
		envManager.Sync(name)
		logger.Separator()
		printValues(name)
	}
}

func process() {
	envManager.SetSystemConfig()
	if listNamesArg {
		listNames()
	} else if variableArg != "" {
		envManager.Sync(variableArg)
		printValues(variableArg)
	} else {
		listFullInfo()
	}
}

func main() {
	cobralib.Run()
}
