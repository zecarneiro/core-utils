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
    if (!$onlyProfile) {
        $dirs = @("$CONFIG_DIR", "$OTHER_APPS_DIR", "$USER_BIN_DIR", "$USER_STARTUP_DIR", "$USER_TEMP_DIR", "$TEMP_DIR")
        Foreach ($dir in $dirs) {
            if (!(Test-Path -Path "$dir")) {
                New-Item -ItemType Directory -Force -Path "$dir" | Out-Null
                Write-Host "Created directory: $dir"
            }
        }
    }
}

function __exit_script {
    exit 0
}

function __set_user_bin_dir {
    $pathKey = "Path"
    $pathEnvArr = ([Environment]::GetEnvironmentVariable($pathKey, [System.EnvironmentVariableTarget]::User) -split ';')
    if (!("$USER_BIN_DIR" -in $pathEnvArr)) {
        $pathEnvArr += "$USER_BIN_DIR"
        [Environment]::SetEnvironmentVariable($pathKey, ($pathEnvArr -join ";"), [System.EnvironmentVariableTarget]::User)
        infolog "Please, Restart the Terminal to change take effect!"
    }
}

function __create_profile_file_powershell {
    if ((fileexists "$MY_SHELL_PROFILE")) {
        Move-Item "$MY_SHELL_PROFILE" "${MY_SHELL_PROFILE}.bk" -Force
    } else {
        infolog "Creating Powershell Script profile to run when powrshell start: $MY_SHELL_PROFILE"
        New-Item "$MY_SHELL_PROFILE" -ItemType file -Force | Out-Null
    }
    if (!(filecontain "$MY_SHELL_PROFILE" "$MY_CUSTOM_SHELL_PROFILE")) {
        writefile "$MY_SHELL_PROFILE" ". '$MY_CUSTOM_SHELL_PROFILE'" -append
    }
    if (!(fileexists "$MY_CUSTOM_SHELL_PROFILE")) {
        infolog "Creating Powershell Script profile to run when powrshell start: $MY_CUSTOM_SHELL_PROFILE"
        New-Item "$MY_CUSTOM_SHELL_PROFILE" -ItemType file -Force | Out-Null
    } 
}

function __install_profile_scripts {
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

function __enable_sudo {
    $res = Read-Host "For Windows 11 only. Do you want to enable sudo? [y/N]"
    if ("$res" -eq "y" -or "$res" -eq "Y") {
        powershell -Command "Start-Process -Wait PowerShell -Verb RunAs -ArgumentList 'sudo.exe config --enable enable'"
    }
}

function __define_default_system_dir {
    $result=$(Read-Host "Insert all User Dirs? (y/N)")
    if ($result -eq "y") {
        $isDocumentsChange = $false
        $currentShellProfileDir = $(dirname "$MY_SHELL_PROFILE")
        $newShellProfileDir = ""
        $userDirs = @{}
        $isSetDirs = $false
        $result=$(selectfolderdialog "Insert DOWNLOAD (Or cancel)")
        if (! [string]::IsNullOrEmpty($result)) {
            $userDirs.Add("{374DE290-123F-4565-9164-39C4925E467B}", "$result")
        }
        $result=$(selectfolderdialog "Insert DOCUMENTS (Or cancel)")
        if (! [string]::IsNullOrEmpty($result)) {
            $userDirs.Add("Personal", "$result")
            $newShellProfileDir = "$result\$(basename "$currentShellProfileDir")"
            $isDocumentsChange = $true
        }
        $result=$(selectfolderdialog "Insert MUSIC (Or cancel)")
        if (! [string]::IsNullOrEmpty($result)) {
            $userDirs.Add("My Music", "$result")
        }
        $result=$(selectfolderdialog "Insert PICTURES (Or cancel)")
        if (! [string]::IsNullOrEmpty($result)) {
            $userDirs.Add("My Pictures", "$result")
        }
        $result=$(selectfolderdialog "Insert VIDEOS (Or cancel)")
        if (! [string]::IsNullOrEmpty($result)) {
            $userDirs.Add("My Video", "$result")
        }
        foreach ($userDir in $userDirs.GetEnumerator()) {
            $isSetDirs=$true
            evaladvanced "reg add `"HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders`" /f /v `"$($userDir.Name)`" /t REG_SZ /d `"$($userDir.Value)`""
        }
        if ($isDocumentsChange) {
            if ((directoryexists "$currentShellProfileDir")) {
                cpdir "$currentShellProfileDir" "$newShellProfileDir"
            }
        }
        if ($isSetDirs){
            restartexplorer
        }
    }
}

function __config_all {
    infolog "No config to process"
}
