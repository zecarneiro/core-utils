package main

import (
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/obj"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "json-log [JSON_DATA]",
		Short: "Print a JSON log message",
		Args:  cobra.MinimumNArgs(1),
	}
	cobralib.WithRunArgsStr(process)
}

func process(data string) {
	var dataObj any
	logic.ProcessError(obj.StringToObject(data, dataObj))
	newData, err := obj.ObjectToString(dataObj)
	logic.ProcessError(err)
	logger.Log(newData)
}

func main() {
	cobralib.Run()
}
