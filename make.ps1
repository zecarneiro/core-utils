$release = $false
$start = $false
$onlyProfile = $false

$SCRIPT_UTILS_DIR = ($PSScriptRoot)
$VERSION = (Get-Content "$SCRIPT_UTILS_DIR\version")
$SHELL_SCRIPT_DIR = "${SCRIPT_UTILS_DIR}\scripts"
$LIBS_DIR = "${SCRIPT_UTILS_DIR}\libs"
$BIN_DIR = "${SCRIPT_UTILS_DIR}\bin"

if ($args[0] -eq "-r" -or $args[0] -eq "--release") {
    $release=$true
}
if ($args[0] -eq "-s" -or $args[0] -eq "--start") {
    $start = $true
}
if ($args[1] -eq "--only-profile-shell") {
    $onlyProfile = $true
}

# ---------------------------------------------------------------------------- #
#                                    IMPORTS                                   #
# ---------------------------------------------------------------------------- #
Get-ChildItem ("${SHELL_SCRIPT_DIR}\*.ps1") | ForEach-Object {
    $fullname = $_.FullName
    . "$fullname"
}
Get-ChildItem ("${LIBS_DIR}\*.ps1") | ForEach-Object {
    $fullname = $_.FullName
    . "$fullname"
}

function process-release {
    . import-libs
    $releasePackageName = "$SCRIPT_UTILS_DIR\core-utils-${VERSION}.zip"
    $releaseDir = "$SCRIPT_UTILS_DIR\release"

    infolog "Create release package"
    evaladvanced "mkdir '$releaseDir'"
    evaladvanced "cpdir '$SHELL_SCRIPT_DIR' '$releaseDir'"
    evaladvanced "cpdir '$LIBS_DIR' '$releaseDir'"
    evaladvanced "cpdir '$BIN_DIR' '$releaseDir'"
    evaladvanced "Copy-Item -Path '$SCRIPT_UTILS_DIR\make.ps1' -Destination '$releaseDir'"
    evaladvanced "Copy-Item -Path '$SCRIPT_UTILS_DIR\make.sh' -Destination '$releaseDir'"
    evaladvanced "Compress-Archive '$releaseDir\*' -DestinationPath '$releasePackageName' -Force"
}

function usage {
    Write-Host "Usage: make.ps1 [OPTIONS]... [STEP-VALUE]"
    Write-Host "OPTIONS:
     -r|--release`tCreate release package
     -s|--start`tProcess install and config by user
    "
}

function printMenu {
    Write-Host "1. Will
    - Enable Sudo
    - Set user bin dir
    - Install Scoop and Winget
    - Install PowershellGet Module
2. Will
    - Install Scoop and Winget packages
    - Install Powershell Modules
    - Install Visual-C-Runtimes
    - Install PIP/PIPX(Depends python language)
    - Install NPM(Include Javascript by default)
    - Install Features for WSL
3. Will
    - Create user powershell profile file
    - Install scripts profile
4. Will
    - Start all configurations for Scoop and Winget
    - Change Full Name of user. User will decide wich to install
    - Install PIPX packages
    - Install Development packages. User will decide wich to install
---
5. Exit"
}

function initProcess {
    $message = "Please, restart your terminal."

    __create_dirs
    if (!(is_valid_home_dir)) {
        show_rules_username
        __exit_script
    }
    for (;;) {
        $option = -1
        if (!($onlyProfile)) {
            printMenu
            $option = Read-Host "Insert an option"
        } else {
            $option = 3
        }
        switch ($option) {
            1 {
                enable-sudo
                set-user-bin-dir
                install-scoop
                install-winget
                install-powershellget
                warnlog "$message"
                __exit_script
            }
            2 {
                install-scoop-packages
                install-winget-packages
                install-modules
                install-visual-c-runtimes
                install-pip-pipx
                install-npm
                __install_features_for_wsl
                warnlog "After reboot, continue with option 3"
                reboot
            }
            3 {
                create-profile-file-powershell
                install-profile-scripts
                if (!$onlyProfile) {
                    Set-Location "$SCRIPT_UTILS_DIR"
                    bash -c "./make.sh --start --only-profile-shell"
                    warnlog "$message"
                }
                __exit_script
            }
            4 {
                config-all
                change_user_full_name
                install-pipx-packages
                install-development-package
                reboot
            }
            5 { __exit_script }
            Default { Write-Host "WARN: Please, insert a valid option!" }
        }
    }
    
}

function main {
    if ($release) {
        process-release
    } elseif ($start) {
        initProcess
    } else {
        usage
    }
}
main
