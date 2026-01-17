package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"golangutils/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	filepath string
	content  string
	mode     string
	encoding string
	forceDir bool
	rootCmd  *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:   "writefile",
		Short: "Write file",
	}
	rootCmd.Flags().StringVarP(&filepath, "file", "f", "", "File to write content")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().StringVarP(&content, "content", "c", "", "Content to write on the file")
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "w", "Mode: 'w' for write, 'a' for append")
	rootCmd.Flags().BoolVarP(&forceDir, "force-dir", "d", false, "Create directory of file if not exists")
	rootCmd.Flags().StringVarP(&encoding, "encoding", "e", "", "Encoding for file. Default: utf-8")
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		var err error
		if !common.InArray([]string{"a", "w"}, mode) {
			logger.Error(fmt.Errorf("invalid mode: %q (must be 'w' or 'a')", mode))
			os.Exit(1)
		}
		isAppend := common.Ternary(mode == "a", true, false)
		if len(encoding) > 0 {
			err = file.WriteFileWithEncoding(filepath, content, isAppend, forceDir, encoding)
		} else {
			err = file.WriteFile(filepath, content, isAppend, forceDir)
		}
		generic.ProcessError(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
