$release = $false
$installStep = 0

$VERSION = "1.0.4"
$SCRIPT_UTILS_DIR = ($PSScriptRoot)
$SHELL_SCRIPT_DIR = "${SCRIPT_UTILS_DIR}\scripts"
$LIBS_DIR = "${SCRIPT_UTILS_DIR}\libs"
$BIN_DIR = "${SCRIPT_UTILS_DIR}\bin"

if ($args[0] -eq "-r" -or $args[0] -eq "--release") {
    $release=$true
}
if ($args[0] -eq "-i" -or $args[0] -eq "--install-step") {
    $installStep = $args[1]
}

function process-release {
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
     -i|--install-step`tProcess install by given step to process
     `tStep 1: Will
     `t`t- Set user bin dir
     `t`t- Install scoop and winget
     `tStep 2: Will
     `t`t- Install scoop and winget packages
     `t`t- Start all configurations for scoop and winget
     `tStep 3: Will
     `t`t- Create user powershell profile file
     `t`t- Install scripts profile
     `tStep 4: Will
     `t`t- Install Visual-C-Runtimes
     `t`t- Install Development packages. User will decide wich to install
    "
}

if ($release -or $installStep -gt 0) {
    Write-Host "INFO: Load scripts..."
    Get-ChildItem ("${SHELL_SCRIPT_DIR}\*.ps1") | ForEach-Object {
        $fullname = $_.FullName
        Write-Host "Loading: $fullname"
        . "$fullname"
    }

    infolog "Load libs..."
    Get-ChildItem ("${LIBS_DIR}\*.ps1") | ForEach-Object {
        $fullname = $_.FullName
        Write-Host "Loading: $fullname"
        . "$fullname"
    }
}

function main {
    $message = "Please, restart your terminal."
    if ($release) {
        process-release
    } elseif ($installStep -eq 1) {
        set-user-bin-dir
        install-scoop
        install-winget
        warnlog "$message"
    } elseif ($installStep -eq 2) {
        install-scoop-packages
        config-all
        warnlog "$message"
    } elseif ($installStep -eq 3) {
        create-profile-file-powershell
        install-profile-scripts
        warnlog "$message"
    } elseif ($installStep -eq 4) {
        install-visual-c-runtimes
        install-development-package
        warnlog "$message"
    } else {
        usage
    }
}
main
