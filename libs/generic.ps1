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

function __create_dirs {
    $dirs = @("$CONFIG_DIR", "$OTHER_APPS_DIR", "$USER_BIN_DIR", "$USER_STARTUP_DIR", "$USER_TEMP_DIR", "$TEMP_DIR")
    Foreach ($dir in $dirs) {
        if (!(Test-Path -Path "$dir")) {
            New-Item -ItemType Directory -Force -Path "$dir" | Out-Null
            Write-Host "Created directory: $dir"
        }
    }
}

function __exit_script {
    exit 0
}

function set-user-bin-dir {
    $pathKey = "Path"
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
    if (!$onlyProfile) {
        infolog "Install core-utils scripts release package"
        Remove-Item -Recurse -Force "$shellScriptsInstallDir" -ErrorAction SilentlyContinue
        cpdir "$SHELL_SCRIPT_DIR" "$shellScriptsInstallDir"
    }
    # Add powershell profiles
    Get-ChildItem -Path "$shellScriptsInstallDir" -Filter *.ps1 -Recurse -File | ForEach-Object {
        $fullName = $_.FullName
        $data = "Import-Module '$fullName' -ErrorAction SilentlyContinue"
        if (!(filecontain "$MY_CUSTOM_SHELL_PROFILE" "$data")) {
            writefile "$MY_CUSTOM_SHELL_PROFILE" "$data" -append
        }
    }
}

function enable-sudo {
    $res = Read-Host "For Windows 11 only. Do you want to enable sudo? [y/N]"
    if ("$res" -eq "y" -or "$res" -eq "Y") {
        powershell -Command "Start-Process -Wait PowerShell -Verb RunAs -ArgumentList 'sudo.exe config --enable enable'"
    }
}

function config-all {
    config-scoop-all
}
