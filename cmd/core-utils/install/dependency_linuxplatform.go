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

func (d *DependencyLinux) installApt(packagesStatus int) {
	cmdInfo := getCmdInfo()
	switch packagesStatus {
	case 0:
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
	case 1:
		logger.Header(fmt.Sprintf("Install APT packages for status: %d", packagesStatus))
		cmdList := []string{
			"sudo apt install zenity -y",
			"sudo apt install dos2unix -y",
			"sudo apt install git -y",
			"sudo apt install powershell -y",
		}
		for _, cmd := range cmdList {
			cmdInfo.Cmd = cmd
			logic.ProcessError(exe.ExecRealTime(cmdInfo))
		}
		envManager.Sync(envPathName)
		logic.ProcessError(exe.ExecRealTime(models.Command{Cmd: "Install-Module -Name PSFzf -Scope CurrentUser", Verbose: true, UseShell: true, ShellToUse: enums.PowerShell}))
	}
}

func (d *DependencyLinux) instalSnap(packagesStatus int) {
	cmdInfo := getCmdInfo()
	if packagesStatus == 0 {
		logger.Header("Install Snap")
		cmdInfo = resetCmdInfo(cmdInfo)
		cmdInfo.Cmd = "sudo apt install snapd-xdg-open snapd -y"
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	}
}

func (d *DependencyLinux) instalFlatpak(packagesStatus int) {
	cmdInfo := getCmdInfo()
	switch packagesStatus {
	case 0:
		logger.Header("Install Flatpak")
		d.runAptCmdWithUpdate("sudo add-apt-repository ppa:flatpak/stable -y")
		cmdList := []string{
			"sudo apt install gnome-software gnome-software-plugin-flatpak xdg-desktop-portal-gtk flatpak -y",
			"flatpak remote-add --if-not-exists flathub 'https://flathub.org/repo/flathub.flatpakrepo'",
		}
		for _, cmd := range cmdList {
			cmdInfo.Cmd = cmd
			logic.ProcessError(exe.ExecRealTime(cmdInfo))
		}
	case 1:
		logger.Header(fmt.Sprintf("Install Flatpak packages for status: %d", packagesStatus))
		libs.RunCoreUtilsCmd("flatpak-install", false, "it.mijorus.gearlever")
	}
}

func (d *DependencyLinux) instalAppImage(packagesStatus int) {
	cmdInfo := getCmdInfo()
	switch packagesStatus {
	case 0:
		logger.Header("Enable AppImage Support")
		cmdList := []string{
			"sudo apt install libfuse2 -y",
			"sudo apt install libfuse2t64 -y",
		}
		for _, cmd := range cmdList {
			cmdInfo.Cmd = cmd
			logic.ProcessError(exe.ExecRealTime(cmdInfo))
		}
	case 1:
		logger.Separator()
		logger.Info("Install mdview")
		mdviewFilepath := file.JoinPath(system.TempDir(), "mdview.AppImage")
		logic.ProcessError(netc.Download("https://github.com/c3er/mdview/releases/download/v3.2.0/mdview-3.2.0-x86_64.AppImage", mdviewFilepath))
		appImageManagerScript := libs.GetScriptAppPathByName("appimage-manager-cu")
		logic.ProcessError(exe.Chmod777(appImageManagerScript, false))
		logic.ProcessError(exe.ExecRealTime(models.Command{
			Cmd:      fmt.Sprintf("%s --install %s", appImageManagerScript, str.GetInSingleQuotes(mdviewFilepath)),
			UseShell: true,
		}))
		logger.Separator()
	}
}

func (d *DependencyLinux) instalPacstall(packagesStatus int) {
	cmdInfo := getCmdInfo()
	switch packagesStatus {
	case 0:
		logger.Header("Install Pacstall")
		cmdInfo.Cmd = "sudo bash -c \"$(curl -fsSL https://pacstall.dev/q/install)\" <<< $'n\n'"
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	case 1:
		logger.Header(fmt.Sprintf("Install Pacstall packages for status: %d", packagesStatus))
		cmdList := []string{
			"sudo pacstall -P -I fzf-bin",
		}
		for _, cmd := range cmdList {
			cmdInfo.Cmd = cmd
			logic.ProcessError(exe.ExecRealTime(cmdInfo))
		}
	}
}

func (d *DependencyLinux) instalDebGet(packagesStatus int) {
	cmdInfo := getCmdInfo()
	switch packagesStatus {
	case 0:
		logger.Header("Install Deb-Get")
		cmdInfo.Cmd = "curl -sL https://raw.githubusercontent.com/wimpysworld/deb-get/main/deb-get | sudo -E bash -s install deb-get"
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
	case 1:
		logger.Header(fmt.Sprintf("Install Deb-Get packages for status: %d", packagesStatus))
		cmdList := []string{
			"sudo deb-get install topgrade",
		}
		for _, cmd := range cmdList {
			cmdInfo.Cmd = cmd
			logic.ProcessError(exe.ExecRealTime(cmdInfo))
		}
	}
}

func (d *DependencyLinux) instalOthers(packagesStatus int) {
	cmdInfo := getCmdInfo()
	switch packagesStatus {
	case 1:
		logger.Header("Install Git Credentials Manager(GCM)")
		tempDirGCM := system.GenerateTempFile("gcm")
		file.CreateDirectory(tempDirGCM, true)
		cmdInfo.Cmd = "sudo bash -c \"$(curl -fsSL https://aka.ms/gcm/linux-install-source.sh)\""
		cmdInfo.Cwd = tempDirGCM
		logic.ProcessError(exe.ExecRealTime(cmdInfo))
		if file.IsDir(tempDirGCM) {
			exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf(`sudo rm -rf "%s"`, tempDirGCM), Verbose: true, UseShell: true})
		}
	}
}

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

func (d *DependencyLinux) start() {
	if console.Confirm("Do you want to install all packages and package managers?", true) {
		if askProcessPackage("Install APT") {
			d.installApt(0)
		}
		if askProcessPackage("Install SNAP") {
			d.instalSnap(0)
		}
		if askProcessPackage("Install FLATPAK") {
			d.instalFlatpak(0)
		}
		if askProcessPackage("Install APPIMAGE") {
			d.instalAppImage(0)
		}
		if askProcessPackage("Install PACKSTALL") {
			d.instalPacstall(0)
		}
		if askProcessPackage("Install DEB-GET") {
			d.instalDebGet(0)
		}
		// Install packages
		envManager.Sync(envPathName)
		if askProcessPackage("Install APT Packages") {
			d.installApt(1)
		}
		if askProcessPackage("Install SNAP Packages") {
			d.instalSnap(1)
		}
		if askProcessPackage("Install FLATPAK Packages") {
			d.instalFlatpak(1)
		}
		if askProcessPackage("Install PACKSTALL Packages") {
			d.instalPacstall(1)
		}
		if askProcessPackage("Install DEB-GET Packages") {
			d.instalDebGet(1)
		}
		if askProcessPackage("Install Others Packages") {
			d.instalOthers(1)
		}
	}
	addUserBinOnPathEnv()
	addCoreUtilsDirsOnPathEnv()
	d.instalScriptsAppsAndAlias()
	if askProcessPackage("Install APPIMAGE Packages") {
		d.instalAppImage(1)
	}
	// Set configs
	d.setConfigs()
	cleanEnvPath()
}
