$release = $false
$start = $false

$VERSION = "1.0.5"
$SCRIPT_UTILS_DIR = ($PSScriptRoot)
$SHELL_SCRIPT_DIR = "${SCRIPT_UTILS_DIR}\scripts"
$LIBS_DIR = "${SCRIPT_UTILS_DIR}\libs"
$BIN_DIR = "${SCRIPT_UTILS_DIR}\bin"

if ($args[0] -eq "-r" -or $args[0] -eq "--release") {
    $release=$true
}
if ($args[0] -eq "-s" -or $args[0] -eq "--start") {
    $start = $true
}

function import-libs {
    Write-Host "INFO: Load scripts..."
    Get-ChildItem ("${SHELL_SCRIPT_DIR}\*.ps1") | ForEach-Object {
        $fullname = $_.FullName
        . "$fullname"
    }

    infolog "Load libs..."
    Get-ChildItem ("${LIBS_DIR}\*.ps1") | ForEach-Object {
        $fullname = $_.FullName
        . "$fullname"
    }
}

function exitScript {
    exit 0
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
    - Install scoop and winget
    - Install PowershellGet Module
2. Will
    - Install scoop and winget packages
    - Install Powershell Modules
    - Start all configurations for scoop and winget
3. Will
    - Create user powershell profile file
    - Install scripts profile
4. Will
    - Install Visual-C-Runtimes
    - Install Development packages. User will decide wich to install
---
5. Exit"
}

function initProcess {
    $message = "Please, restart your terminal."
    for (;;) {
        printMenu
        $option = Read-Host "Insert an option"
        if ($option -gt 0 -and $option -lt 5) {
            . import-libs
        }
        switch ($option) {
            1 {
                enable-sudo
                set-user-bin-dir
                install-scoop
                install-winget
                install-powershellget
                warnlog "$message"
                exitScript
            }
            2 {
                install-scoop-packages
                install-modules
                config-all
                warnlog "$message"
                exitScript
            }
            3 {
                create-profile-file-powershell
                install-profile-scripts
                warnlog "$message"
                exitScript
            }
            4 {
                install-visual-c-runtimes
                install-development-package
                warnlog "$message"
                exitScript
            }
            5 { exitScript }
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
