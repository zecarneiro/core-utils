package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/env"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"main/internal/dir"
	"main/internal/libs"
	"main/internal/libs/cobralib"
	"regexp"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

const (
	BRAVE_BROWSER_NAME         = "brave"
	GOOGLE_CHROME_BROWSER_NAME = "google-chrome"
	MS_EDGE_BROWSER_NAME       = "microsoft-edge"
)

var (
	nameArg           string
	urlArg            string
	browserArg        string
	browserBinArg     string
	startMaximizedArg bool
	iconArg           string
	isolateDataArg    bool
	bashCmdArg        string

	validBrowserList []string
	browserBin       string
	validIconExt     string
	webAppsDataDir   string
	appDataDir       string
	nameNormalized   string
	cmdFilePath      string
	iconDstDir       string
)

func init() {
	console.EnableFeatures()
	logger.WithSeparatorLength(10)
	loadData()
	setupCommand()
}

func loadData() {
	browserBin = ""
	validIconExt = logic.Ternary(platform.IsWindows(), "ico", "png")
	webAppsDataDir = fmt.Sprintf(`%s/web-apps-data`, system.HomeUserLocalDir())
	validBrowserList = []string{BRAVE_BROWSER_NAME, GOOGLE_CHROME_BROWSER_NAME, MS_EDGE_BROWSER_NAME}
}

func setupCommand() {
	cobralib.CobraCmd = &cobra.Command{
		Use:   "webapp-installer",
		Short: "Install Web Apps by given browser, url and Name. Only work with Chrommium browser based. Will be create a command, shurtcut and folder for data located on system local dir.",
	}
	cobralib.CobraCmd.Flags().StringVarP(&nameArg, "name", "n", "", "Name of App")
	cobralib.CobraCmd.Flags().StringVarP(&urlArg, "url", "u", "", "URL of App")
	cobralib.CobraCmd.Flags().StringVarP(&browserArg, "browser", "b", "", fmt.Sprintf(`Indicate a browser to use. Accept only: %s`, slice.ArrayToStringBySep(validBrowserList, ",")))
	cobralib.CobraCmd.Flags().StringVarP(&browserBinArg, "browser-executable", "e", "", "Command to run installed Browser. If set, browser arg will be ignored.")
	cobralib.CobraCmd.Flags().BoolVarP(&startMaximizedArg, "start-maximized", "m", false, "App window will be started maximized.")
	cobralib.CobraCmd.Flags().StringVarP(&iconArg, "icon", "i", "", "Icon of App")
	cobralib.CobraCmd.Flags().BoolVarP(&isolateDataArg, "isolate-data", "D", false, "Isolate App data dir, means will be created a separeted user data dir.")
	if !platform.IsWindows() {
		cobralib.CobraCmd.Flags().StringVarP(&bashCmdArg, "bash-command", "B", console.WhichIgnoreError("bash"), "Bash command to run created app command.")
	}
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("name"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("url"))
	logic.ProcessError(cobralib.CobraCmd.MarkFlagRequired("browser"))
	cobralib.WithRun(process)
}

func getRunner() string {
	if platform.IsWindows() {
		runnerDir := fmt.Sprintf(`%s\web-apps-runner`, dir.CoreUtilsLocal())
		if err := file.CreateDirectory(runnerDir, true); err != nil {
			logger.Error(err)
		} else {
			output := libs.RunCoreUtilsCmdWithOutput("create-vbs-runner", false, "-n", nameNormalized, "-c", cmdFilePath)
			if !str.IsEmpty(output) {
				runnerSrc := strings.TrimSpace(str.StringReplaceAll(output, map[string]string{"Created:": ""}))
				runnerDst := fmt.Sprintf(`%s\%s`, runnerDir, file.Basename(runnerSrc))
				if err := file.Move(runnerSrc, runnerDir); err != nil {
					logger.Error(err)
				} else {
					return runnerDst
				}
			}
		}
	}
	return cmdFilePath
}

func getIcon() string {
	if file.IsDir(iconDstDir) && file.IsFile(iconArg) {
		var err error
		iconDst := file.ResolvePath(fmt.Sprintf(`%s/%s.%s`, iconDstDir, nameNormalized, validIconExt))
		logger.Separator()
		if file.IsFileExtension(iconArg, validIconExt) {
			logger.Info(fmt.Sprintf(`Copy icon from: %s`, iconArg))
			logger.Info(fmt.Sprintf(`To: %s`, iconDst))
			err = file.CopyFile(iconArg, file.ResolvePath(iconDst))
		} else {
			logger.Info(fmt.Sprintf(`Convert icon from: %s`, iconArg))
			logger.Info(fmt.Sprintf(`To: %s`, iconDst))
			err = exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf(`magick "%s" "%s"`, iconArg, iconDst), UseShell: true, Verbose: false})
		}
		if err == nil {
			return iconDst
		}
		logger.Error(err)
	}
	return ""
}

func createMenu() {
	createMenuEntryCmd := "create-menu-entry"
	icon := getIcon()
	runner := getRunner()
	cmdArgs := []string{
		"-n", nameArg,
		"-e", logic.Ternary(platform.IsWindows(), runner, fmt.Sprintf(`%s -c '%s'`, bashCmdArg, runner)),
	}
	if !str.IsEmpty(icon) {
		cmdArgs = append(cmdArgs, "-i", icon)
	}
	logger.Separator()
	logger.Info(fmt.Sprintf(`Will be use on menu entry this runner/command: %s`, runner))
	libs.RunCoreUtilsCmd(createMenuEntryCmd, false, cmdArgs...)
}

func getCommand() string {
	cmdData := fmt.Sprintf(`%s --user-data-dir='%s'`, browserBin, appDataDir)
	if startMaximizedArg {
		cmdData = fmt.Sprintf(`%s --start-maximized`, cmdData)
	}
	cmdData = fmt.Sprintf(`%s --app='%s'`, cmdData, urlArg)
	return cmdData
}

func validate() {
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	if !slices.Contains(validBrowserList, browserArg) {
		logic.ProcessError(fmt.Errorf("Insert a valid browser. It's important to define browser data dir. If diferent browser have different data dir, the app can have error."))
	}
	if str.IsEmpty(nameArg) || str.IsEmpty(nameNormalized) || !re.MatchString(nameNormalized) {
		logic.ProcessError(fmt.Errorf("Invalid given name: %s. Name must not be empty and accept only A-Z, a-z, 0-9, - and _", nameArg))
	}
	if err := file.CreateDirectory(appDataDir, true); err != nil {
		logic.ProcessError(err)
	}
	if str.IsEmpty(browserBin) {
		logic.ProcessError(fmt.Errorf(`Not found bin of given browser: %s`, browserArg))
	}
	if !platform.IsWindows() && !file.IsFile(bashCmdArg) {
		logic.ProcessError(fmt.Errorf(`Can not find the bash cmd. Please, set this argument.`))
	}
	if !str.IsEmpty(iconArg) {
		err := file.CreateDirectory(iconDstDir, true)
		if err != nil {
			logger.ErrorStr(fmt.Sprintf("Error on create icon directory: %s", iconDstDir))
			logger.Error(err)
		}
	}

}

func fillAllNecessaryVar() {
	nameNormalized = strings.ToLower(strings.ReplaceAll(nameArg, " ", "-"))
	appDataDir = file.ResolvePath(fmt.Sprintf(`%s/%s`, webAppsDataDir, logic.Ternary(isolateDataArg, nameNormalized, browserArg)))
	if isolateDataArg {
		iconDstDir = appDataDir
	} else {
		iconDstDir = file.ResolvePath(fmt.Sprintf(`%s/web-app-installer-icons`, appDataDir))
	}
	browserBin = browserBinArg
	if platform.IsWindows() {
		cmdFilePath = fmt.Sprintf(`%s\%s.cmd`, dir.CoreUtilsPrompt(), nameNormalized)
	} else {
		cmdFilePath = fmt.Sprintf(`%s/%s`, dir.CoreUtilsShellScripts(), nameNormalized)
	}
	if str.IsEmpty(browserBin) {
		findFirstBrowserExec := func(exePaths []string) string {
			exePathFounded := ""
			for _, exePath := range exePaths {
				if file.IsFile(exePath) {
					exePathFounded = exePath
					break
				}
			}
			return exePathFounded
		}
		switch browserArg {
		case BRAVE_BROWSER_NAME:
			if platform.IsWindows() {
				browserBin = findFirstBrowserExec([]string{
					`C:\Program Files\BraveSoftware\Brave-Browser\Application\brave.exe`,
					fmt.Sprintf(`%s\BraveSoftware\Brave-Browser\Application\brave.exe`, env.GetWithSingleValue("LOCALAPPDATA")),
				})
			} else if platform.IsLinux() {
				browserBin = console.WhichIgnoreError("brave-browser")
			}
		case GOOGLE_CHROME_BROWSER_NAME:
			if platform.IsWindows() {
				browserBin = findFirstBrowserExec([]string{
					`C:\Program Files\Google\Chrome\Application\chrome.exe`,
					fmt.Sprintf(`%s\Google\Chrome\Application\chrome.exe`, env.GetWithSingleValue("LOCALAPPDATA")),
				})
			} else if platform.IsLinux() {
				browserBin = console.WhichIgnoreError("google-chrome")
			}
		case MS_EDGE_BROWSER_NAME:
			if platform.IsWindows() {
				browserBin = findFirstBrowserExec([]string{
					`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
					fmt.Sprintf(`%s\Microsoft\Edge\Application\msedge.exe`, env.GetWithSingleValue("LOCALAPPDATA")),
				})
			} else if platform.IsLinux() {
				browserBin = console.WhichIgnoreError("microsoft-edge")
			}
		}
	}
}

func process() {
	logger.Info(fmt.Sprintf(`Installing WebApp: %s`, nameArg))
	fillAllNecessaryVar()
	validate()
	libs.RunCoreUtilsCmd("alias-manager-cu", false, "-n", nameNormalized, "-f", "-c", getCommand())
	createMenu()
	logger.Separator()
	logger.Info("To manage created commands, use command: script-manager-cu")
	logger.Info(fmt.Sprintf(`This app will use this user data dir: %s`, appDataDir))
	logger.Ok("Installation is done.")
}

func main() {
	cobralib.Run()
}
