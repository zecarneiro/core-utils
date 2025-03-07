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

function __show_install_message_question {
    param([string] $message)
    $userInput = (Read-Host "Do you want to install ${message}? (y/N)")
    return $userInput
}

function set-user-bin-dir {
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

function create-profile-file-powershell {
    if ((fileexists "$MY_SHELL_PROFILE")) {
        Move-Item "$MY_SHELL_PROFILE" "${MY_SHELL_PROFILE}.bk" -Force
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

function install-profile-scripts {
    $shellScriptsInstallDir = "${OTHER_APPS_DIR}\shell-scripts"
    infolog "Install core-utils scripts release package"

    Remove-Item -Recurse -Force "$shellScriptsInstallDir" -ErrorAction SilentlyContinue
    cpdir "$SHELL_SCRIPT_DIR" "${OTHER_APPS_DIR}"
    Rename-Item "${OTHER_APPS_DIR}\scripts" "$shellScriptsInstallDir"
    # Add powershell profiles
    Get-ChildItem -Path "$shellScriptsInstallDir" -Filter *.ps1 -Recurse -File | ForEach-Object {
        $fullName = $_.FullName
        $data = "Import-Module '$fullName' -ErrorAction SilentlyContinue"
        if (!(filecontain "$MY_CUSTOM_SHELL_PROFILE" "$data")) {
            writefile "$MY_CUSTOM_SHELL_PROFILE" "$data" -append
        }
    }
}

function config-all {
    config-scoop-all
    $res = Read-Host "For Windows 11 only. Do you want to enable sudo? [y/N]"
    if ("$res" -eq "y" -or "$res" -eq "Y") {
        powershell -Command "Start-Process -Wait PowerShell -Verb RunAs -ArgumentList 'sudo.exe config --enable enable'"
    }
}
