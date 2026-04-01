package main

import (
	"errors"
	"fmt"
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
	appName        = "service-manager-cu"
	nameArg        string
	commandArg     string
	commandArgsArg string
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
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "Service name (required)")
	cobralib.CobraCmd.Flags().StringVarP(&commandArg, "command", "c", "", "Command to run")
	cobralib.CobraCmd.Flags().StringVarP(&commandArgsArg, "command-args", "a", "", "Arguments for the command")
	cobralib.CobraCmd.Flags().StringVarP(&descriptionArg, "description", "d", "", "Service description")
	cobralib.CobraCmd.Flags().BoolVarP(&installArg, "install", "i", false, "Install service")
	cobralib.CobraCmd.Flags().BoolVarP(&uninstallArg, "uninstall", "u", false, "Uninstall service")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("name"))
	cobralib.CobraCmd.MarkFlagsMutuallyExclusive("install", "uninstall")
	cobralib.CobraCmd.MarkFlagsOneRequired("install", "uninstall")
	cobralib.WithRun(process)
}

func getServiceDir() string {
	serviceDir := ""
	if platform.IsWindows() {
		serviceDir = file.ResolvePath("C:\\services")
	} else if platform.IsLinux() {
		serviceDir = file.ResolvePath("/etc/systemd/system")
	}
	if !file.IsDir(serviceDir) {
		logic.ProcessError(file.CreateDirectory(serviceDir, true))
	}
	return serviceDir
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
	commandExt := file.FileExtension(commandArg)
	if platform.IsWindows() && commandExt != "exe" {
		logic.ProcessError(fmt.Errorf("On Windows, all services must be an .exe file. This command extension is: %s. Use task manager instead.", commandExt))
	}
}

func process() {
	if !platform.IsWindows() && !platform.IsLinux() {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
	validateName()
	validateCommand()
	serviceDir := getServiceDir()
	if uninstallArg {
		uninstall(serviceDir)
	} else if installArg {
		install(serviceDir)
	} else {
		logic.ProcessError(cobralib.CobraCmd.Help())
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

func install(serviceDir string) {
	if platform.IsWindows() {
		serviceCmd := fmt.Sprintf(`%s %s`, commandArg, commandArgsArg)
		cmds := []string{
			fmt.Sprintf("sc.exe create \"%s\" binPath=\"%s\" type=own start=auto", nameArg, serviceCmd),
			fmt.Sprintf("sc.exe start \"%s\"", nameArg),
		}
		for _, cmd := range cmds {
			runCommand(cmd, true)
		}
	} else if platform.IsLinux() {
		serviceCmd := fmt.Sprintf(`%s %s`, commandArg, commandArgsArg)
		serviceFile := file.JoinPath(serviceDir, fmt.Sprintf("%s.service", nameArg))
		content := fmt.Sprintf(linuxTemplate, descriptionArg, serviceCmd)
		fileWriter := models.FileWriterConfig{
			File:        serviceFile,
			Data:        content,
			IsAppend:    false,
			WithUtf8BOM: false,
			IsCreateDir: true,
		}
		logic.ProcessError(file.WriteFile(fileWriter))
		cmds := []string{
			fmt.Sprintf("systemctl start %s", nameArg),
			fmt.Sprintf("systemctl enable %s", nameArg),
		}
		for _, cmd := range cmds {
			runCommand(cmd, true)
		}
	}
	logger.Ok(fmt.Sprintf("Installed service with name: %s", nameArg))
}

func uninstall(serviceDir string) {
	if platform.IsWindows() {
		cmds := []string{
			fmt.Sprintf("sc.exe stop \"%s\"", nameArg),
			fmt.Sprintf("sc.exe delete \"%s\"", nameArg),
		}
		for _, cmd := range cmds {
			runCommand(cmd, false)
		}
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
