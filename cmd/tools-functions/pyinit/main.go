package main

import (
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "pyinit",
		Short: "Create an empty __init__.py file in the current directory",
	}
	cobralib.WithRun(process)
}

func process() {
	cwd, err := file.GetCurrentDir()
	logic.ProcessError(err)
	initFilePath := file.JoinPath(cwd, "__init__.py")
	fileConfig := models.FileWriterConfig{
		File:        initFilePath,
		Data:        "",
		IsAppend:    false,
		IsCreateDir: false,
	}
	logic.ProcessError(file.WriteFile(fileConfig))
}

func main() {
	cobralib.Run()
}
