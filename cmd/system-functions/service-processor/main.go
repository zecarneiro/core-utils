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
	cobralib.CobraCmd.Flags().StringVarP(&descriptionArg, "description", "d", "", "Service description")
	cobralib.CobraCmd.Flags().BoolVarP(&installArg, "install", "i", false, "Install service")
	cobralib.CobraCmd.Flags().BoolVarP(&uninstallArg, "uninstall", "u", false, "Uninstall service")
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("name"))
	cobralib.WithRun(process)
}

func getServiceDir() string {
	serviceDir := ""
	if platform.IsWindows() {
		serviceDir = file.ResolvePath(file.JoinPath("C:", "services"))
	} else if platform.IsLinux() {
		serviceDir = file.ResolvePath("/etc/systemd/system")
	}
	logic.ProcessError(errors.New(platform.UnsupportedMSG))
	if !file.IsDir(serviceDir) {
		logic.ProcessError(file.CreateDirectory(serviceDir, true))
	}
	return serviceDir
}

func validateName(name string) {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if str.IsEmpty(name) || !re.MatchString(name) {
		logic.ProcessError(fmt.Errorf("%s: Invalid given name: %s. Name must not be empty and accept only A-Z, a-z, 0-9, - and _", appName, name))
	}
}

func process() {
	if !platform.IsWindows() && !platform.IsLinux() {
		logic.ProcessError(errors.New(platform.UnsupportedMSG))
	}
	validateName(nameArg)
	serviceDir := getServiceDir()
	if installArg && str.IsEmpty(commandArg) {
		logic.ProcessError(fmt.Errorf("invalid given command: %s", commandArg))
	}
	if uninstallArg {
		uninstall(serviceDir)
	} else if installArg {
		install(serviceDir)
	} else {
		logic.ProcessError(cobralib.CobraCmd.Help())
	}
}

func install(serviceDir string) {
	if platform.IsWindows() {
		serviceFile := file.JoinPath(serviceDir, fmt.Sprintf("%s.bat", nameArg))
		content := fmt.Sprintf(windowsTemplate, descriptionArg, commandArg)
		fileWriter := models.FileWriterConfig{
			File:        serviceFile,
			Data:        content,
			IsAppend:    false,
			WithUtf8BOM: false,
			IsCreateDir: true,
		}
		logic.ProcessError(file.WriteFile(fileWriter))
		cmds := []string{
			fmt.Sprintf("sc.exe create \"%s\" binPath=\"%s\" type=own start=auto", nameArg, serviceFile),
			fmt.Sprintf("sc.exe start \"%s\"", nameArg),
		}
		for _, cmd := range cmds {
			execCmd := models.Command{
				Cmd:      cmd,
				UseShell: true,
				Verbose:  true,
			}
			logic.ProcessError(exe.ExecRealTime(execCmd))
		}
	} else if platform.IsLinux() {
		serviceFile := file.JoinPath(serviceDir, fmt.Sprintf("%s.service", nameArg))
		content := fmt.Sprintf(linuxTemplate, descriptionArg, commandArg, nameArg, nameArg)
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
			execCmd := models.Command{
				Cmd:      cmd,
				Verbose:  true,
				UseShell: true,
			}
			logic.ProcessError(exe.ExecRealTime(execCmd))
		}
	}
	logger.Ok(fmt.Sprintf("Installed service with name: %s", nameArg))
}

func uninstall(serviceDir string) {
	if platform.IsWindows() {
		serviceFile := file.JoinPath(serviceDir, fmt.Sprintf("%s.bat", nameArg))
		cmds := []string{
			fmt.Sprintf("sc.exe stop \"%s\"", nameArg),
			fmt.Sprintf("sc.exe delete \"%s\"", nameArg),
		}
		for _, cmd := range cmds {
			execCmd := models.Command{
				Cmd:      cmd,
				UseShell: true,
				Verbose:  true,
			}
			logic.ProcessError(exe.ExecRealTime(execCmd))
		}
		logic.ProcessError(file.DeleteFile(serviceFile))
	} else if platform.IsLinux() {
		serviceFile := file.JoinPath(serviceDir, fmt.Sprintf("%s.service", nameArg))
		cmds := []string{
			fmt.Sprintf("systemctl stop %s", nameArg),
		}
		for _, cmd := range cmds {
			execCmd := models.Command{
				Cmd:     cmd,
				Verbose: true,
			}
			logic.ProcessError(exe.ExecRealTime(execCmd))
		}
		logic.ProcessError(file.DeleteFile(serviceFile))
	}
	logger.Ok(fmt.Sprintf("Uninstaled service with name: %s", nameArg))
}

func main() {
	cobralib.Run()
}
