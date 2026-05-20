package main

import (
	"fmt"
	"golangutils/pkg/system"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var isFullName bool

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "OS_NAME",
		Short: "Name of Operating System",
	}
	cobralib.CobraCmd.Flags().BoolVarP(&isFullName, "full", "f", false, "Print full name of Operating System.")
	cobralib.WithRun(process)
}

func process() {
	name := strings.ToUpper(string(system.GetOsType().String()))
	if isFullName {
		name = system.OSFullName()
	}
	fmt.Println(name)
}

func main() {
	cobralib.Run()
}
