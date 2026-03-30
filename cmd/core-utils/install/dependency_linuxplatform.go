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
	"golangutils/pkg/netc"
	"golangutils/pkg/shell"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"golangutils/pkg/time"

	"main/internal/dir"
	"main/internal/libs"
)

type DependencyLinux struct{}

func NewDependencyLinux() *DependencyLinux {
	return &DependencyLinux{}
}

/* -------------------------------------------------------------------------- */
/*                                     APT                                    */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) runAptUpdate() {
	cmdInfo := getCmdInfo()
	cmdInfo.Cmd = "sudo apt update"
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
}

func (d *DependencyLinux) runAptCmdWithUpdate(cmd string) {
	cmdInfo := getCmdInfo()
	cmdInfo.Cmd = cmd
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
	d.runAptUpdate()
}

func (d *DependencyLinux) installApt() {
	cmdInfo := getCmdInfo()
	logger.Header("Install APT and Importants packages")
	d.runAptUpdate()
	d.runAptCmdWithUpdate("sudo add-apt-repository universe -y")
	d.runAptCmdWithUpdate("sudo add-apt-repository multiverse -y")
	cmdList := []string{
		"sudo apt install apt-transport-https -y",
		"sudo apt install wget -y",
		"sudo apt install curl -y",
		"sudo apt install libnotify-bin -y",
		"sudo apt install lsb-release -y",
	}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = cmd
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	}
	logger.Separator()
	logger.Info("Install microsoft repository")
	microsoftRepoDebFile := "/tmp/packages-microsoft-prod.deb"
	downloadUrl := fmt.Sprintf("https://packages.microsoft.com/config/ubuntu/%s/packages-microsoft-prod.deb", system.OSVersion())
	logic.ProcessError(netc.Download(downloadUrl, microsoftRepoDebFile))
	cmdInfo = resetCmdInfo(cmdInfo)
	d.runAptCmdWithUpdate(fmt.Sprintf("sudo dpkg -i '%s'", microsoftRepoDebFile))
	time.Sleep(5)
	logic.ProcessError(file.DeleteFile(microsoftRepoDebFile))
	logger.Separator()
}

func (d *DependencyLinux) installAptPackages() {
	cmdInfo := getCmdInfo()
	logger.Header("Install APT packages")
	cmdList := []string{"zenity", "dos2unix", "git", "powershell"}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = fmt.Sprintf(`sudo apt install %s -y`, cmd)
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
	envManager.Sync(envPathName)
	logger.Error(exe.ExecRealTime(models.Command{Cmd: "Install-Module -Name PSFzf -Scope CurrentUser -SkipPublisherCheck", Verbose: true, UseShell: true, ShellToUse: enums.PowerShell}))
}

/* -------------------------------------------------------------------------- */
/*                                    SNAP                                    */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) instalSnap() {
	cmdInfo := getCmdInfo()
	logger.Header("Install Snap")
	cmdList := []string{"snapd-xdg-open", "snapd"}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = fmt.Sprintf(`sudo apt install %s -y`, cmd)
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
}

/* -------------------------------------------------------------------------- */
/*                                   FLATPAK                                  */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) instalFlatpak() {
	cmdInfo := getCmdInfo()
	logger.Header("Install Flatpak")
	d.runAptCmdWithUpdate("sudo add-apt-repository ppa:flatpak/stable -y")
	cmdAptList := []string{"gnome-software", "gnome-software-plugin-flatpak", "xdg-desktop-portal-gtk", "flatpak"}
	for _, cmd := range cmdAptList {
		cmdInfo.Cmd = fmt.Sprintf(`sudo apt install %s -y`, cmd)
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
	cmdList := []string{
		"flatpak remote-add --if-not-exists flathub 'https://flathub.org/repo/flathub.flatpakrepo'",
	}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = cmd
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	}
}

func (d *DependencyLinux) instalFlatpakPackages() {
	logger.Header("Install Flatpak packages")
	libs.RunCoreUtilsCmd("flatpak-install", false, "it.mijorus.gearlever")
}

/* -------------------------------------------------------------------------- */
/*                                  APPIMAGE                                  */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) instalAppImage() {
	cmdInfo := getCmdInfo()
	logger.Header("Enable AppImage Support")
	cmdAptList := []string{"libfuse2", "libfuse2t64"}
	for _, cmd := range cmdAptList {
		cmdInfo.Cmd = fmt.Sprintf(`sudo apt install %s -y`, cmd)
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
}

func (d *DependencyLinux) instalAppImagePackages() {
	logger.Header("Install AppImage packages")
	logger.Separator()
	logger.Info("Install mdview")
	mdviewFilepath := file.JoinPath(system.TempDir(), "mdview.AppImage")
	logger.Error(netc.Download("https://github.com/c3er/mdview/releases/download/v3.2.0/mdview-3.2.0-x86_64.AppImage", mdviewFilepath))
	appImageManagerScript := libs.GetScriptAppPathByName("appimage-manager-cu")
	logger.Error(exe.Chmod777(appImageManagerScript, false))
	logger.Error(exe.ExecRealTime(models.Command{
		Cmd:      fmt.Sprintf("%s --install %s", appImageManagerScript, str.GetInSingleQuotes(mdviewFilepath)),
		UseShell: true,
	}))
	logger.Separator()
}

/* -------------------------------------------------------------------------- */
/*                                  PACKSTALL                                 */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) instalPacstall() {
	cmdInfo := getCmdInfo()
	logger.Header("Install Pacstall")
	cmdInfo.Cmd = "sudo bash -c \"$(curl -fsSL https://pacstall.dev/q/install)\" <<< $'n\n'"
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
}

func (d *DependencyLinux) instalPacstallPackages() {
	cmdInfo := getCmdInfo()
	logger.Header("Install Pacstall packages")
	cmdList := []string{"fzf-bin"}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = fmt.Sprintf(`sudo pacstall -P -I %s`, cmd)
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
}

/* -------------------------------------------------------------------------- */
/*                                   DEB-GET                                  */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) instalDebGet() {
	cmdInfo := getCmdInfo()
	logger.Header("Install Deb-Get")
	cmdInfo.Cmd = "curl -sL https://raw.githubusercontent.com/wimpysworld/deb-get/main/deb-get | sudo -E bash -s install deb-get"
	logic.ProcessError(exe.ExecRealTime(cmdInfo))
}

func (d *DependencyLinux) instalDebGetPackages() {
	cmdInfo := getCmdInfo()
	logger.Header("Install Deb-Get packages")
	cmdList := []string{
		"topgrade",
		"gcm", // git-credential-manager
	}
	for _, cmd := range cmdList {
		cmdInfo.Cmd = fmt.Sprintf(`sudo deb-get install %s`, cmd)
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
}

/* -------------------------------------------------------------------------- */
/*                                   OTHERS                                   */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) instalScriptsAppsAndAlias() {
	envManager.Sync(envPathName)
	scriptAppsDir := file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "apps", "bash")
	filesInfo, err := file.ReadDirRecursive(scriptAppsDir)
	logic.ProcessError(err)
	// Install coreutils profile scripts
	installShellScriptOnSystemProfile(shell.GetShellProfileFile(enums.Bash), fmt.Sprintf("source '%s'", file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "profiles", "init-bash-shell.sh")))
	installShellScriptOnSystemProfile(shell.GetShellProfileFile(enums.PowerShell), fmt.Sprintf(". \"%s\"", file.JoinPath(dir.CoreUtilsSystemInstallShellScripts(), "profiles", "init-pwsh-shell.ps1")))
	// Install All apps
	logger.Header("Install all apps")
	for _, fileInfo := range filesInfo.Files {
		libs.RunCoreUtilsCmd("script-manager-cu", false, "install", fileInfo)
	}
	// Add alias
	logger.Header("Install all alias")
	alias := map[string]string{
		"restart-pipewire":    "systemctl --user restart pipewire.service",
		"pwsh":                fmt.Sprintf("%s -nologo", shell.GetPowershellCmd()),
		"powershell":          "pwsh",
		"zenity":              fmt.Sprintf("%s 2>/dev/null", console.WhichIgnoreError("zenity")),
		"now":                 "date",
		"sha1":                "openssl sha1",
		"md5":                 "openssl md5",
		"sha256":              "openssl sha256",
		"pause":               "echo -n \"Press [ENTER] to continue...: \"; read var_name",
		"cls":                 "clear",
		"update-menu-entries": "sudo update-desktop-database",
		"gearlever":           "flatpak run it.mijorus.gearlever",
	}
	for key, value := range alias {
		libs.RunCoreUtilsCmd("alias-manager-cu", false, "-n", key, "-c", value)
	}
}

func (d *DependencyLinux) setConfigs() {
	// Set config again, because the last time I ran, the powershell is not installed
	envManager.SetSystemConfig()
	envManager.Sync(envPathName)
	// Create menu entries
	powershellMenuEntryArgs := []string{
		"-n", "Powershell",
		"-e", shell.GetPowershellCmd(),
		"-a", "-nologo",
		"-i", "utilities-terminal",
		"-c", "ConsoleOnly;System;",
		"-t",
	}
	libs.RunCoreUtilsCmd("create-menu-entry", false, powershellMenuEntryArgs...)
}

/* -------------------------------------------------------------------------- */
/*                                    START                                   */
/* -------------------------------------------------------------------------- */
func (d *DependencyLinux) start() {
	if console.Confirm("Do you want to install all packages and package managers?", true) {
		// PACKAGE APP
		if askProcessPackage("Install APT") {
			d.installApt()
		}
		if askProcessPackage("Install SNAP") {
			d.instalSnap()
		}
		if askProcessPackage("Install FLATPAK") {
			d.instalFlatpak()
		}
		if askProcessPackage("Install APPIMAGE") {
			d.instalAppImage()
		}
		if askProcessPackage("Install PACKSTALL") {
			d.instalPacstall()
		}
		if askProcessPackage("Install DEB-GET") {
			d.instalDebGet()
		}
		// Update ENV's
		envManager.Sync(envPathName)
		// INSTALL APP PACKAGES
		if askProcessPackage("Install APT Packages") {
			d.installAptPackages()
		}
		if askProcessPackage("Install FLATPAK Packages") {
			d.instalFlatpakPackages()
		}
		if askProcessPackage("Install PACKSTALL Packages") {
			d.instalPacstallPackages()
		}
		if askProcessPackage("Install DEB-GET Packages") {
			d.instalDebGetPackages()
		}
	}
	addUserBinOnPathEnv()
	addCoreUtilsDirsOnPathEnv()
	d.instalScriptsAppsAndAlias()
	if askProcessPackage("Install APPIMAGE Packages") {
		d.instalAppImagePackages()
	}
	// Set configs
	d.setConfigs()
	cleanEnvPath()
}
