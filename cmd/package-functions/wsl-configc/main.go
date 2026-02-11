package main

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/system"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var ram, processor int

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "wsl-configc",
		Short: "This configurations only works on windows 11 or newer!!",
	}
	cobralib.CobraCmd.Flags().IntVarP(&ram, "ram", "r", -1, "Max RAM(GB) that WSL will use")
	cobralib.CobraCmd.Flags().IntVarP(&processor, "processor", "p", -1, "Max Processor that WSL will use")
	cobralib.WithRun(process)
}

func process() {
	configFile := file.ResolvePath(system.HomeDir(), ".wslconfig")
	data := ""
	if file.IsFile(configFile) {
		backup := fmt.Sprintf("%s.bk", configFile)
		logger.Warn(fmt.Sprintf("Config file already exists: %s", configFile))
		logger.Warn(fmt.Sprintf("Backup original config file to: %s", backup))
		logic.ProcessError(file.CopyFile(configFile, backup))
		logic.ProcessError(file.DeleteFile(configFile))
	}
	if ram > 0 {
		data = fmt.Sprintf("memory=%dGB", ram)
	}
	if processor > 0 {
		data = logic.Ternary(len(data) > 0, fmt.Sprintf("%s%sprocessors=%d", data, common.Eol(), processor), fmt.Sprintf("processors=%d", processor))
	}
	if len(data) > 0 {
		data = fmt.Sprintf("[wsl2]%s%s", common.Eol(), data)
		logic.ProcessError(file.WriteFile(configFile, data, false, false))
		// TODO: wsl_shutdown()
	}
}

func main() {
	cobralib.Run()
}
