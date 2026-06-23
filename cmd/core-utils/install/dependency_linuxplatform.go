package main

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/conv"
	"golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/shell"
	"golangutils/pkg/system"

	"main/internal/dir"
	"main/internal/libs"
)

type DependencyLinux struct {
	isUpdate bool
}

func NewDependencyLinux(isUpdate bool) *DependencyLinux {
	return &DependencyLinux{isUpdate: isUpdate}
}

/* -------------------------------------------------------------------------- */
/*                                   OTHERS                                   */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) instalScriptsAppsAndAlias() {
	scriptAppsDir := file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "apps", "bash")
	// Install coreutils profile scripts
	installShellScriptOnSystemProfile(shell.GetShellProfileFile(enums.Bash), fmt.Sprintf("source '%s'", file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "profiles", "init-bash-shell.sh")))
	installShellScriptOnSystemProfile(shell.GetShellProfileFile(enums.PowerShell), fmt.Sprintf(". \"%s\"", file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "profiles", "init-pwsh-shell.ps1")))
	// Install All apps
	logger.Header("Install all apps")
	runScript("install-scripts-apps-and-alias.sh", "scripts-apps", conv.BoolToString(d.isUpdate), scriptAppsDir)
	// Add alias
	logger.Header("Install all alias")
	runScript("install-scripts-apps-and-alias.sh", "alias", conv.BoolToString(d.isUpdate))
}

func (d *DependencyLinux) createMenus() {
	logger.Header("Create Menu Entries")
	runScript("install-scripts-apps-and-alias.sh", "menu-entries")
}

func (d *DependencyLinux) setConfigs() {
	// Set config again, because the last time I ran, the powershell is not installed
	envManager.SetSystemConfig()
	envManager.Sync(envPathName)
	d.createMenus()
}

func (d *DependencyLinux) configDataPartition() {
	logger.Header("Configure partitions to auto mount")
	autoMountScript := libs.GetScriptAppPathByName("auto-mount-ntfs-ext4")
	logger.Error(exe.Chmod777(autoMountScript, false))
	logger.Error(exe.ExecRealTime(models.Command{
		Cmd:      autoMountScript,
		UseShell: true,
	}))
}

/* -------------------------------------------------------------------------- */
/*                                START/UPDATE                                */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) update() {
	d.instalScriptsAppsAndAlias()
	d.createMenus()
}

func (d *DependencyLinux) install() {
	if console.Confirm("Do you want to install all packages and package managers?", true) {
		// PACKAGE APP
		if askProcessPackage("Install APT and Importants packages") {
			runScript("install-apt-and-packages.sh", "manager", system.OSVersion())
		}
		if askProcessPackage("Install SNAP") {
			runScript("install-snap-and-packages.sh", "manager")
		}
		if askProcessPackage("Install FLATPAK") {
			runScript("install-flatpak-and-packages.sh", "manager")
		}
		if askProcessPackage("Install APPIMAGE") {
			runScript("install-appimage-and-packages.sh", "manager")
		}
		if askProcessPackage("Install PACKSTALL") {
			runScript("install-pacstall-and-packages.sh", "manager")
		}
		if askProcessPackage("Install DEB-GET") {
			runScript("install-deb-get-and-packages.sh", "manager")
		}
		// INSTALL APP PACKAGES
		if askProcessPackage("Install APT Packages") {
			runScript("install-apt-and-packages.sh", "packages")
		}
		if askProcessPackage("Install FLATPAK Packages") {
			runScript("install-flatpak-and-packages.sh", "packages")
		}
		if askProcessPackage("Install PACKSTALL Packages") {
			runScript("install-pacstall-and-packages.sh", "packages")
		}
		if askProcessPackage("Install github Packages") {
			runScript("install-github-packages.sh", "packages")
		}
		if askProcessPackage("Install APPIMAGE Packages") {
			runScript("install-appimage-and-packages.sh", "packages")
		}
	}
	addUserBinOnPathEnv()
	addCoreUtilsDirsOnPathEnv()
	d.instalScriptsAppsAndAlias()
	libs.RunCoreUtilsCmd("installer-by-url", false, "run")

	// Set configs
	d.setConfigs()
	cleanEnvPath()
	d.configDataPartition()
	changeAndCreateDefaultDirs()
}

func (d *DependencyLinux) start() {
	if d.isUpdate {
		d.update()
	} else {
		d.install()
	}
}
