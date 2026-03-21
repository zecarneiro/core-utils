package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"slices"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	filepath string
	content  string
	mode     string
	encoding string
	forceDir bool
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "writefile",
		Short: "Write file",
	}
	cobralib.CobraCmd.Flags().StringVarP(&filepath, "file", "f", "", "File to write content")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("file"))
	cobralib.CobraCmd.Flags().StringVarP(&content, "content", "c", "", "Content to write on the file")
	cobralib.CobraCmd.Flags().StringVarP(&mode, "mode", "m", "w", "Mode: 'w' for write, 'a' for append")
	cobralib.CobraCmd.Flags().BoolVarP(&forceDir, "force-dir", "d", false, "Create directory of file if not exists")
	cobralib.CobraCmd.Flags().StringVarP(&encoding, "encoding", "e", "", "Encoding for file. Default: utf-8")
	cobralib.WithRun(process)
}

func process() {
	if !slices.Contains([]string{"a", "w"}, mode) {
		logic.ProcessError(fmt.Errorf("invalid mode: %q (must be 'w' or 'a')", mode))
	}
	fileConfig := models.FileWriterConfig{
		File:         filepath,
		Data:         content,
		IsAppend:     logic.Ternary(mode == "a", true, false),
		IsCreateDir:  forceDir,
		EncodingName: encoding,
	}
	logic.ProcessError(file.WriteFile(fileConfig))
}

func main() {
	cobralib.Run()
}
