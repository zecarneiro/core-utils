package main

import (
	"errors"
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"regexp"

	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	appName        = "task-manager-cu"
	nameArg        string
	commandArg     string
	commandArgs    []string
	descriptionArg string
	installArg     bool
	uninstallArg   bool
)

func init() {
	if !system.IsAdmin() {
		logic.ProcessError(errors.New(system.NeedAdminAccessMsg))
	}
	setupCommand()
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   appName,
		Short: "Install or uninstall a service",
	}
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "Task name (required)")
	cobralib.CobraCmd.Flags().StringArrayVarP(&commandArgs, "command", "c", []string{""}, "Command to run")
	cobralib.CobraCmd.Flags().StringVarP(&descriptionArg, "description", "d", "", "Task description")
	cobralib.CobraCmd.Flags().BoolVarP(&installArg, "install", "i", false, "Install task")
	cobralib.CobraCmd.Flags().BoolVarP(&uninstallArg, "uninstall", "u", false, "Uninstall task")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("name"))
	cobralib.CobraCmd.MarkFlagsMutuallyExclusive("install", "uninstall")
	cobralib.CobraCmd.MarkFlagsOneRequired("install", "uninstall")
	cobralib.WithRun(process)
}

func getTaskDir() string {
	taskDir := ""
	if platform.IsWindows() {
		taskDir = file.ResolvePath("C:\\services")
	} else if platform.IsLinux() {
		taskDir = file.JoinPath("/etc/systemd/system")
	}
	if !file.IsDir(taskDir) {
		logic.ProcessError(file.CreateDirectory(taskDir, true))
	}
	return taskDir
}

func mapCommand() {
	commandArg = ""
	for _, cmdArg := range commandArgs {
		if !str.IsEmpty(cmdArg) {
			if platform.IsWindows() {
				commandArg = logic.Ternary(str.IsEmpty(commandArg), cmdArg, fmt.Sprintf("%s & %s", commandArg, cmdArg))
			} else if platform.IsLinux() {
				commandArg += fmt.Sprintf("ExecStart=%s%s", cmdArg, common.Eol())
			}
		}
	}
}

func validateName() {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if str.IsEmpty(nameArg) || !re.MatchString(nameArg) {
		logic.ProcessError(fmt.Errorf("%s: Invalid given name: %s. Name must not be empty and accept only A-Z, a-z, 0-9, - and _", appName, nameArg))
	}
}

func validateCommand() {
	if installArg && str.IsEmpty(commandArg) {
		logic.ProcessError(fmt.Errorf("invalid given command: %s", commandArg))
	}
}

func runCommand(cmd string, isThrown bool) {
	execCmd := models.Command{
		Cmd:      cmd,
		UseShell: true,
		Verbose:  true,
	}
	if isThrown {
		logic.ProcessError(exe.ExecRealTime(execCmd))
	} else {
		logger.Error(exe.ExecRealTime(execCmd))
	}
}

func process() {
	if !platform.IsWindows() && !platform.IsLinux() {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
	mapCommand()
	validateName()
	validateCommand()
	taskDir := getTaskDir()
	if uninstallArg {
		uninstall(taskDir)
	} else if installArg {
		install(taskDir)
	} else {
		logic.ProcessError(cobralib.CobraCmd.Help())
	}
}

func install(taskDir string) {
	if platform.IsWindows() {
		taskFile := file.JoinPath(taskDir, fmt.Sprintf("%s.bat", nameArg))
		content := fmt.Sprintf(windowsTemplate, commandArg)
		fileWriter := models.FileWriterConfig{
			File:        taskFile,
			Data:        content,
			IsAppend:    false,
			WithUtf8BOM: false,
			IsCreateDir: true,
		}
		logic.ProcessError(file.WriteFile(fileWriter))
		cmds := []string{
			fmt.Sprintf(`schtasks /create /tn "%s" /sc onlogon /tr "cmd.exe /c %s" /f`, nameArg, taskFile),
			fmt.Sprintf("schtasks /run /tn \"%s\"", nameArg),
		}
		for _, cmd := range cmds {
			runCommand(cmd, true)
		}
	} else if platform.IsLinux() {
		serviceFile := file.JoinPath(taskDir, fmt.Sprintf("%s.service", nameArg))
		content := fmt.Sprintf(linuxTemplate, descriptionArg, commandArg)
		fileWriter := models.FileWriterConfig{
			File:        serviceFile,
			Data:        content,
			IsAppend:    false,
			WithUtf8BOM: false,
			IsCreateDir: true,
		}
		logic.ProcessError(file.WriteFile(fileWriter))
		cmds := []string{
			"systemctl daemon-reload",
			fmt.Sprintf("systemctl start %s.service", nameArg),
			fmt.Sprintf("systemctl enable %s.service", nameArg),
		}
		for _, cmd := range cmds {
			runCommand(cmd, true)
		}
	}
	logger.Ok(fmt.Sprintf("Installed service with name: %s", nameArg))
}

func uninstall(serviceDir string) {
	if platform.IsWindows() {
		serviceFile := file.JoinPath(serviceDir, fmt.Sprintf("%s.bat", nameArg))
		cmds := []string{
			fmt.Sprintf("schtasks /end /tn \"%s\"", nameArg),
			fmt.Sprintf("schtasks /delete /tn \"%s\" /f", nameArg),
		}
		for _, cmd := range cmds {
			runCommand(cmd, false)
		}
		logic.ProcessError(file.DeleteFile(serviceFile))
	} else if platform.IsLinux() {
		serviceFile := file.JoinPath(serviceDir, fmt.Sprintf("%s.service", nameArg))
		cmds := []string{
			fmt.Sprintf("systemctl stop %s", nameArg),
			fmt.Sprintf("systemctl disable %s", nameArg),
		}
		for _, cmd := range cmds {
			runCommand(cmd, false)
		}
		logic.ProcessError(file.DeleteFile(serviceFile))
		runCommand("systemctl daemon-reload", false)
		runCommand("systemctl reset-failed", false)
	}
	logger.Ok(fmt.Sprintf("Uninstaled service with name: %s", nameArg))
}

func main() {
	cobralib.Run()
}
