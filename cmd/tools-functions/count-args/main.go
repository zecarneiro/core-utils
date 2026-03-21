package main

import (
	"fmt"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "count-args",
		Short: "Count the number of arguments passed",
	}
	cobralib.WithRun(process)
}

func process() {
	args := cobralib.GetArgs()
	fmt.Println(len(args))
}

func main() {
	cobralib.Run()
}
