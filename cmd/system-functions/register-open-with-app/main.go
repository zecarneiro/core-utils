package main

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"strings"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	friendlyName string
	isUninstall  bool
)

func init() { setupCommand() }

func setOthersArg(cmdToSet *cobra.Command) {
	cmdToSet.Flags().StringVarP(&friendlyName, "friendly-name", "f", "", "Set friendly name fot the app (required, Not if delete flag is not set))")
	cmdToSet.Flags().BoolVarP(&isUninstall, "delete", "d", false, "Unregister app)")
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "register-open-with-app",
		Short: "Register app into open with. IMPORTANT: On the end of the process, will restart the explorer.",
	}

	// Scoop
	scoopCmd := &cobra.Command{
		Use:   "scoop",
		Short: "Register scoop app",
		Run: func(cmd *cobra.Command, args []string) {
			scoopProcess(cmd)
		},
	}
	scoopCmd.Flags().StringP("app-name", "n", "", "Set app name")
	scoopCmd.MarkFlagRequired("app-name")
	setOthersArg(scoopCmd)

	// Scoop
	executableCmd := &cobra.Command{
		Use:   "executable",
		Short: "Register executable app",
		Run: func(cmd *cobra.Command, args []string) {
			executableProcess(cmd)
		},
	}
	executableCmd.Flags().StringP("path", "p", "", "Set full path for executable")
	executableCmd.MarkFlagRequired("path")
	setOthersArg(executableCmd)
	cobralib.CobraCmd.AddCommand(scoopCmd, executableCmd)
}

func getExeRegName(exeFileName string) string {
	return fmt.Sprintf(`cu_%s`, exeFileName)
}

func getRegFile(appName string) string {
	return file.JoinPath(system.TempDir(), fmt.Sprintf(`%s.reg`, appName))
}

func run(regTempFile string) {
	logger.Error(exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf("sudo cmd.exe /C %s", regTempFile), Verbose: true, UseShell: true}))
	logger.Error(file.DeleteFile(regTempFile))
	libs.RunCoreUtilsCmd("restart-explorer", true)
}

func startRegisterApp(appName string, exeFileName string, exePath string) {
	regTempFile := getRegFile(appName)
	exeRegName := getExeRegName(exeFileName)
	templateData := fmt.Sprintf(REG_ADD_TEMPLATE, exeRegName, friendlyName, exeRegName, exeRegName, exeRegName, strings.ReplaceAll(exePath, `\`, `\\`), "%1")
	logic.ProcessError(file.WriteFile(models.FileWriterConfig{File: regTempFile, Data: templateData, IsAppend: false, IsCreateDir: true}))
	run(regTempFile)
}

func validate() {
	if !isUninstall {
		if str.IsEmpty(friendlyName) {
			logic.ProcessError(fmt.Errorf("Invalid given friendly name."))
		}
	}
}

func scoopProcess(cmd *cobra.Command) {
	validate()
	appName, err := cmd.Flags().GetString("app-name")
	logic.ProcessError(err)
	if str.IsEmpty(appName) {
		logic.ProcessError(fmt.Errorf("Invalid given app name."))
	}
	_, exeFileName, exePath := getScoopAppInfo(appName)
	if isUninstall {
		uninstallProcess(appName, exeFileName)
	} else {
		startRegisterApp(appName, exeFileName, exePath)
	}
}

func executableProcess(cmd *cobra.Command) {
	validate()
	executable, err := cmd.Flags().GetString("path")
	logic.ProcessError(err)
	if !file.IsFile(executable) {
		logic.ProcessError(fmt.Errorf("Invalid given executable."))
	}
	exeFileName := file.Basename(executable)
	appName := file.FileName(exeFileName)
	if isUninstall {
		uninstallProcess(appName, exeFileName)
	} else {
		startRegisterApp(appName, exeFileName, executable)
	}
}

func uninstallProcess(appName string, exeFileName string) {
	validate()
	regTempFile := getRegFile(appName)
	exeRegName := getExeRegName(exeFileName)
	templateData := fmt.Sprintf(REG_DELETE_TEMPLATE, exeRegName)
	logic.ProcessError(file.WriteFile(models.FileWriterConfig{File: regTempFile, Data: templateData, IsAppend: false, IsCreateDir: true}))
	run(regTempFile)
}

func main() {
	cobralib.Run()
}
