package main

import (
	"fmt"
	"strings"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	textArg      string
	delimiterArg string
	directionArg string
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "cutc <text>",
		Short: "Cut text by delimiter and print left or right side",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.CobraCmd.Flags().StringVarP(&delimiterArg, "delimiter", "d", "", "Delimiter (required)")
	cobralib.CobraCmd.Flags().StringVarP(&directionArg, "direction", "dr", "", "Direction: L or R (required)")
	cobralib.CobraCmd.MarkFlagRequired("delimiter")
	cobralib.CobraCmd.MarkFlagRequired("direction")
	cobralib.WithRunArgsStr(process)
}

func process(text string) {
	dataResult := strings.SplitN(text, delimiterArg, 2)
	switch directionArg {
	case "L":
		fmt.Println(dataResult[0])
	case "R":
		if len(dataResult) > 1 {
			fmt.Println(dataResult[1])
		} else {
			fmt.Println("")
		}
	default:
		fmt.Println("")
	}
}

func main() {
	cobralib.Run()
}
