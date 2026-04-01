package main

import (
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"os"
	"regexp"

	"main/internal/libs"
	"main/internal/libs/cobralib"

	"github.com/spf13/cobra"
)

var (
	nameArg    string
	commandArg string
	waitArg    bool

	templateLinux = `[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Hidden=false
NoDisplay=true
X-GNOME-Autostart-enabled=true
`
)

func init() { setupCommand() }

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "create-boot-app-entry",
		Short: "Create a Boot APP for system",
	}
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "VBS filename (required)")
	cobralib.CobraCmd.Flags().StringVarP(&commandArg, "command", "c", "", "Command to run or binary (required)")
	cobralib.WithRun(process)
}

func validateName(name string) {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if str.IsEmpty(name) || !re.MatchString(name) {
		logic.ProcessError(fmt.Errorf("Invalid given name!. Name must not be empty and accept only A-Z, a-z, 0-9, - and _"))
	}
}

func process() {
	validateName(nameArg)
	if str.IsEmpty(commandArg) {
		logic.ProcessError(fmt.Errorf("Command is required"))
	}

	if platform.IsWindows() {
		appDataDir := os.Getenv("APPDATA")
		if str.IsEmpty(appDataDir) {
			logic.ProcessError(fmt.Errorf("Not found APPDATA env"))
		}
		bootDir := fmt.Sprintf(`%s\Microsoft\Windows\Start Menu\Programs\Startup`, appDataDir)
		libs.RunCoreUtilsCmd("create-vbs-runner", false, "-n", nameArg, "-c", commandArg, "-w", bootDir)
	} else {
		bootDir := file.JoinPath(system.HomeUserConfigDir(), "autostart")
		logic.ProcessError(file.CreateDirectory(bootDir, true))
		content := fmt.Sprintf(templateLinux, nameArg, commandArg)
		fileConfig := models.FileWriterConfig{
			File:        fmt.Sprintf("%s.desktop", file.JoinPath(bootDir, nameArg)),
			Data:        content,
			IsAppend:    false,
			WithUtf8BOM: false,
			IsCreateDir: true,
		}
		logic.ProcessError(file.WriteFile(fileConfig))
	}
}

func main() {
	cobralib.Run()
}
