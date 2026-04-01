package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/shell"
	"golangutils/pkg/system"

	"main/internal/dir"
	"main/internal/libs"
)

type DependencyWindows struct {
	scoopDir   string
	gitBashBin string
	isUpdate   bool
}

func NewDependencyWindows(isUpdate bool) *DependencyWindows {
	scoopDir := file.JoinPath(system.HomeDir(), "scoop")
	return &DependencyWindows{
		scoopDir:   scoopDir,
		gitBashBin: file.JoinPath(scoopDir, "apps", "git", "current", "bin", "bash.exe"),
		isUpdate:   isUpdate,
	}
}

/* -------------------------------------------------------------------------- */
/*                                   WINGET                                   */
/* -------------------------------------------------------------------------- */
func (d *DependencyWindows) installWinget() {
	logger.Header("Install Winget")
	wingetInstallerScript := libs.GetScriptCmdPathByName("install-winget.ps1", "core-utils")
	cmd := models.Command{
		Cmd:      "sudo",
		Args:     []string{shell.GetPowershellCmd(), fmt.Sprintf(`-Command "%s"`, wingetInstallerScript)},
		UseShell: true,
		Verbose:  true,
	}
	logger.Error(exe.ExecRealTime(cmd))
}

func (d *DependencyWindows) installWingetPackages() {
	logger.Header("Install Winget packages")
	idList := []string{
		"topgrade-rs.topgrade",
	}
	for _, id := range idList {
		libs.RunCoreUtilsCmd("winget-install", false, id)
	}
}

/* -------------------------------------------------------------------------- */
/*                                    SCOOP                                   */
/* -------------------------------------------------------------------------- */
func (d *DependencyWindows) installScoop() {
	cmdInfo := getCmdInfo()
	logger.Header("Install Scoop")
	cmdList := []string{
		"Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression",
		"Install-Module -AllowClobber -Name scoop-completion -Scope CurrentUser -SkipPublisherCheck", // Project URL - https://github.com/Moeologist/scoop-completion",
	}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = cmd
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	}
}

func (d *DependencyWindows) installScoopPackages() {
	cmdInfo := getCmdInfo()
	logger.Header("Install Scoop packages")
	cmdList := []string{
		"%s install main/7zip",
		"%s install main/git",
		"%s install main/openssl",
		"%s install main/vim",
		"%s install main/nano",
		"%s install main/curl",
		"%s install main/grep",
		"%s install main/sed",
		"%s install main/clink",
		"%s install main/uutils-coreutils",
		"%s install main/dos2unix",
		"%s install main/fzf",
		"%s bucket add extras",
		"%s install extras/psfzf",
		"%s install extras/psreadline", // https://github.com/PowerShell/PSReadLine
		"%s install extras/git-credential-manager",
		"%s install https://github.com/c3er/mdview/releases/latest/download/mdview.json", // https://github.com/c3er/mdview
	}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = fmt.Sprintf(cmd, "scoop")
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
}

func (d *DependencyWindows) configScoopPackages() {
	cmdInfo := getCmdInfo()
	logger.Header("Config Scoop packages")
	cmdList := []string{
		fmt.Sprintf("sudo cmd.exe /C %s\\apps\\7zip\\current\\install-context.reg", d.scoopDir),
		"clink set clink.logo none",
	}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = cmd
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
}

/* -------------------------------------------------------------------------- */
/*                                   OTHERS                                   */
/* -------------------------------------------------------------------------- */
func (d *DependencyWindows) createMenus() {
	logger.Header("Create Menu Entries")
	wtBin, err := console.Which("wt.exe")
	if err == nil {
		bashMenuEntryArgs := []string{
			"-n", "Bash",
			"-e", fmt.Sprintf("`\"%s`\"", wtBin),
			"-a", fmt.Sprintf("--title Bash -d `\"%s`\" `\"%s`\" --login -i", system.HomeDir(), d.gitBashBin),
			"-c", "ConsoleOnly;System;",
			"-i", "C:\\Windows\\System32\\cmd.exe,0",
		}
		libs.RunCoreUtilsCmd("create-menu-entry", false, bashMenuEntryArgs...)
	}
	logger.Error(err)
}

func (d *DependencyWindows) setConfigs() {
	envManager.Sync(envPathName)
	// Create menu entries
	d.createMenus()
}

func (d *DependencyWindows) instalScriptsAppsAndAlias() {
	envManager.Sync(envPathName)
	forceArg := logic.Ternary(d.isUpdate, "-f", "")
	scriptAppsDir := file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "apps", "pwsh")
	filesInfo, err := file.ReadDirRecursive(scriptAppsDir)
	logic.ProcessError(err)
	// Install All apps
	logger.Header("Install all apps")
	for _, fileInfo := range filesInfo.Files {
		libs.RunCoreUtilsCmd("script-manager-cu", false, "install", fileInfo, forceArg)
	}
	// Add alias
	logger.Header("Install all alias")
	alias := map[string]string{
		"pwsh":       fmt.Sprintf("%s -nologo", shell.GetPowershellCmd()),
		"powershell": "pwsh",
		"now":        "date",
		"bash":       d.gitBashBin,
		"pausec":     "powershell.cmd pause",
	}
	for key, value := range alias {
		isShellNoArgs := logic.Ternary(key == "now", "-o", "")
		libs.RunCoreUtilsCmd("alias-manager-cu", false, "-n", key, "-c", value, isShellNoArgs, forceArg)
	}
	logger.Header("Disable system alias")
	systemAlias := []string{
		"cp",
		"cat",
		"mkdir",
		"ls",
		"mv",
		"ps",
		"rm",
		"rmdir",
		"sleep",
		"sort",
		"tee",
		"curl",
		"grep",
		"sed",
	}
	for _, value := range systemAlias {
		libs.RunCoreUtilsCmd("disable-system-alias-cu", false, "-n", value)
	}
	// Install coreutils profile scripts
	installShellScriptOnSystemProfile(shell.GetShellProfileFile(enums.Bash), fmt.Sprintf("source '%s'", file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "profiles", "init-bash-shell.sh")))
	installShellScriptOnSystemProfile(shell.GetShellProfileFile(enums.PowerShell), fmt.Sprintf(". \"%s\"", file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "profiles", "init-pwsh-shell.ps1")))
}

/* -------------------------------------------------------------------------- */
/*                                START/UPDATE                                */
/* -------------------------------------------------------------------------- */
func (d *DependencyWindows) update() {
	d.instalScriptsAppsAndAlias()
	d.createMenus()
}

func (d *DependencyWindows) install() {
	if console.Confirm("Do you want to install all packages and package managers?", true) {
		if askProcessPackage("Install Winget") {
			d.installWinget()
		}
		if askProcessPackage("Install Scoop") {
			d.installScoop()
		}
		if askProcessPackage("Enable WSL") {
			cmdWsl := getCmdInfo()
			cmdWsl.Cmd = "wsl.exe --install"
			logger.Error(exe.ExecRealTime(cmdWsl))
		}
		// Install packages
		envManager.Sync(envPathName)
		if askProcessPackage("Install Winget Packages") {
			d.installWingetPackages()
		}
		if askProcessPackage("Install Scoop Packages") {
			d.installScoopPackages()
		}
		// Config packages
		envManager.Sync(envPathName)
		if askProcessPackage("Config Scoop Packages") {
			d.configScoopPackages()
		}
	}
	addUserBinOnPathEnv()
	addCoreUtilsDirsOnPathEnv()
	d.instalScriptsAppsAndAlias()
	d.setConfigs()
	cleanEnvPath()
	libs.RunCoreUtilsCmd("change-user-full-name", false)
	changeAndCreateDefaultDirs()
}

func (d *DependencyWindows) start() {
	if d.isUpdate {
		d.update()
	} else {
		d.install()
	}
}
