package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/conv"
	"golangutils/pkg/enums"
	"golangutils/pkg/env"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/shell"
	"golangutils/pkg/system"

	"main/internal/dir"
	"main/internal/libs"
)

type DependencyWindows struct {
	scoopDir string
	isUpdate bool
}

func NewDependencyWindows(isUpdate bool) *DependencyWindows {
	scoopDir := file.JoinPath(system.HomeDir(), "scoop")
	return &DependencyWindows{
		scoopDir: scoopDir,
		isUpdate: isUpdate,
	}
}

/* -------------------------------------------------------------------------- */
/*                                   OTHERS                                   */
/* -------------------------------------------------------------------------- */
func (d *DependencyWindows) createMenus() {
	logger.Header("Create Menu Entries")
	runScript("install-scripts-apps-and-alias.ps1", "menu-entries")
}

func (d *DependencyWindows) setConfigs() {
	envManager.SetSystemConfig()
	envManager.Sync(envPathName)
	d.createMenus()
}

func (d *DependencyWindows) instalScriptsAppsAndAlias() {
	scriptAppsDir := file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "apps", "pwsh")
	forceArg := fmt.Sprintf(`$%s`, conv.BoolToString(d.isUpdate))
	// Install All apps
	logger.Header("Install all apps")
	runScript("install-scripts-apps-and-alias.ps1", "scripts-apps", "-FORCE_ARG", forceArg, "-SCRIPTS_APPS_DIR_ARG", scriptAppsDir)
	// Add alias
	logger.Header("Install all alias")
	runScript("install-scripts-apps-and-alias.ps1", "alias", "-FORCE_ARG", forceArg, "-SCOOP_DIR_ARG", d.scoopDir)
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
		coreUtilsExecDir := dir.CoreutilsExecDir(true)
		env.InsertOnPath(coreUtilsExecDir)
		// PACKAGE APP
		if askProcessPackage("Install WINGET") {
			runScript("install-winget-and-packages.ps1", "manager", "-WingetInstallerScript", libs.GetScriptCmdPathByName("install-winget.ps1", "core-utils"))
		}
		if askProcessPackage("Install SCOOP") {
			runScript("install-scoop-and-packages.ps1", "manager")
		}
		if askProcessPackage("Enable WSL") {
			runScript("install-wsl.ps1", "manager")
		}
		// INSTALL APP PACKAGES
		if askProcessPackage("Install Winget Packages") {
			runScript("install-winget-and-packages.ps1", "packages")
		}
		if askProcessPackage("Install Scoop Packages") {
			runScript("install-scoop-and-packages.ps1", "packages")
		}
		// Config packages
		if askProcessPackage("Config Scoop Packages") {
			runScript("install-scoop-and-packages.ps1", "config", "-SCOOP_DIR_ARG", d.scoopDir)
		}
	}
	addUserBinOnPathEnv()
	addCoreUtilsDirsOnPathEnv()
	d.instalScriptsAppsAndAlias()
	libs.RunCoreUtilsCmd("installer-by-url", false, "run")

	// Set configs
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
