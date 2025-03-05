$VERSION = "1.0.1"
$SCRIPT_UTILS_DIR = ($PSScriptRoot)
$SHELL_SCRIPT_DIR = "${SCRIPT_UTILS_DIR}\scripts"
$SHELL_SCRIPT_TEMP_DIR = "$([System.IO.Path]::GetTempPath())core-utils"

function __commandexists($command) {
    $oldPreference = $ErrorActionPreference
    $ErrorActionPreference = 'SilentlyContinue'
    try {
        if (Get-Command $command) {
            RETURN $true
        }
    } Catch {
        RETURN $false
    } Finally {
        $ErrorActionPreference = $oldPreference
    }
}

function release {
    $releasePackageName = "$SCRIPT_UTILS_DIR\core-utils-${VERSION}.zip"  
    Write-Host "INFO: Create release package"
    Compress-Archive "$SHELL_SCRIPT_DIR\*" -DestinationPath "$releasePackageName" -Force
}

function install_scripts {
    $shellScriptsInstallDir = "${OTHER_APPS_DIR}\shell-scripts"
    Write-Host "INFO: Install core-utils scripts release package"

    Remove-Item -Recurse -Force "$shellScriptsInstallDir" -ErrorAction SilentlyContinue
    New-Item -Path "$shellScriptsInstallDir" -ItemType Directory -Force
    # Add powershell profiles
    Get-ChildItem -Path "$SHELL_SCRIPT_TEMP_DIR" -Filter *.ps1 -Recurse -File | ForEach-Object {
        $fullName = $_.FullName
        $name = $_.Name
        $fullInstallName = "$shellScriptsInstallDir\_$name"
        cp "$fullName" "$fullInstallName"
        $data = "Import-Module '$fullInstallName' -ErrorAction SilentlyContinue"
        if (!(filecontain "$MY_CUSTOM_SHELL_PROFILE" "$data")) {
            writefile "$MY_CUSTOM_SHELL_PROFILE" "$data" -append
        }
    }
}

function create_profile_file_powershell {
    if ((fileexists "$MY_SHELL_PROFILE")) {
        mv "$MY_SHELL_PROFILE" "${MY_SHELL_PROFILE}.bk"
    } else {
        infolog "Creating Powershell Script profile to run when powrshell start: $MY_SHELL_PROFILE"
        New-Item "$MY_SHELL_PROFILE" -ItemType file -Force
    }
    if (!(filecontain "$MY_SHELL_PROFILE" "$MY_CUSTOM_SHELL_PROFILE")) {
        writefile "$MY_SHELL_PROFILE" ". '$MY_CUSTOM_SHELL_PROFILE'" -append
    }
    if (!(fileexists "$MY_CUSTOM_SHELL_PROFILE")) {
        infolog "Creating Powershell Script profile to run when powrshell start: $MY_CUSTOM_SHELL_PROFILE"
        New-Item "$MY_CUSTOM_SHELL_PROFILE" -ItemType file -Force
    } 
}

function process_scoop_packages {
    evaladvanced "scoop install main/coreutils"
    evaladvanced "scoop install main/git"
    evaladvanced "scoop install main/vim"
    evaladvanced "scoop install main/nano"
    evaladvanced "scoop install main/curl"
    evaladvanced "scoop install main/grep"
    evaladvanced "scoop install main/sed"
    evaladvanced "scoop install main/which"
    evaladvanced "scoop install main/dos2unix"
    evaladvanced "scoop bucket add extras"
    evaladvanced "scoop install extras/okular"
    delalias "cp"
    delalias "cat"
    delalias "mkdir"
    delalias "ls"
    delalias "mv"
    delalias "ps"
    delalias "rm"
    delalias "rmdir"
    delalias "sleep"
    delalias "sort"
    delalias "tee"
    delalias "curl"
    delalias "grep"
    delalias "sed"
}

function install_scoop {
    if (!(__commandexists scoop)) {
        Write-Host "INFO: Install Scoop ..."
        Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
        Install-Module -AllowClobber -Name scoop-completion -Scope CurrentUser # Project URL - https://github.com/Moeologist/scoop-completion
    }
}

function set_user_bin_dir {
    $pathKey = "Path"
    if (!(Test-Path -Path "$USER_BIN_DIR")) {
        New-Item -ItemType Directory -Path "$USER_BIN_DIR" | Out-Null
    }
    $pathEnvArr = ([Environment]::GetEnvironmentVariable($pathKey, [System.EnvironmentVariableTarget]::User) -split ';')
    if (!("$USER_BIN_DIR" -in $pathEnvArr)) {
        $pathEnvArr += "$USER_BIN_DIR"
        [Environment]::SetEnvironmentVariable($pathKey, ($pathEnvArr -join ";"), [System.EnvironmentVariableTarget]::User)
        infolog "Please, Restart the Terminal to change take effect!"
    }
}

function downloadCoreUtilsScripts {
    $coreUtilPackage = "$([System.IO.Path]::GetTempPath())core-utils.zip"
    $url = "https://github.com/zecarneiro/core-utils/releases/download/v${VERSION}/core-utils-${VERSION}.zip"

    Write-Host "INFO: Clean core-utils scripts release package"
    Remove-Item "$SHELL_SCRIPT_TEMP_DIR" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item "$coreUtilPackage" -Recurse -Force -ErrorAction SilentlyContinue

    Write-Host "INFO: Download core-utils scripts release package"
    Invoke-WebRequest "$url" -OutFile "$coreUtilPackage"
    Expand-Archive -Path "$coreUtilPackage" -DestinationPath "$SHELL_SCRIPT_TEMP_DIR"
}

if ($global:installStep -eq 1) {
    downloadCoreUtilsScripts
}
if ($global:installStep -eq 1 -or $global:installStep -eq 2 -or $global:release) {
    Write-Host "INFO: Load scripts..."
    Get-ChildItem ("${SHELL_SCRIPT_TEMP_DIR}\*.ps1") | ForEach-Object {
        $fullname = $_.FullName
        Write-Host "Loading: $fullname"
        . "$fullname"
    }
}

function main {
    Write-Host "INFO: Received arguments: InstallStep=${global:installStep}, Release=${release}"
    $message = "INFO: Please, restart your terminal."
    if ($global:installStep -eq 1) {
        set_user_bin_dir
        install_scoop
        Write-Host "$message"
    } elseif ($global:installStep -eq 2) {
        . reloadprofile
        process_scoop_packages
        . reloadprofile
        create_profile_file_powershell
        install_scripts
        setautoloadmodule "scoop-completion"
        Write-Host "$message"
    } elseif ($global:release) {
        release
    } else {
        Write-Host "Usage: make.ps1 - Global Vars [installStep|release]=[$true|]"
    }
}
main
