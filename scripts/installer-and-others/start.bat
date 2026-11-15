@echo off
cd /d "%~dp0"

set coreutil_rootdir=%cd%
if "%1"=="first-config" (
    echo [INFO] Enable Sudo
    powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "%coreutil_rootdir%\scripts\apps\pwsh\runadmin.ps1 'sudo.exe config --enable enable'"
    echo [INFO] Enabling Powershell Policy for all users...
    sudo powershell.exe -NoProfile -Command "Set-ExecutionPolicy RemoteSigned -Scope LocalMachine -Force"
    echo [INFO] Exclude Coreutils folder from antivirus
    sudo powershell.exe -NoProfile -Command "%coreutil_rootdir%\scripts\cmd\core-utils\exclude-folder-for-ms-security.ps1"
    echo [INFO] Process some system policies
    sudo powershell.exe -NoProfile -Command "%coreutil_rootdir%\scripts\apps\pwsh\windows-policies-processor.ps1"
    echo [INFO] Restarting the PC to changes take efect...
    sudo powershell.exe -NoProfile -Command "Restart-Computer -Force"
) else if "%1"=="install" (
    echo Installing...
    sudo powershell.exe -NoProfile -Command "Add-MpPreference -ExclusionPath '%coreutil_rootdir%'"
    "%coreutil_rootdir%\installers\install.exe" "%coreutil_rootdir%"
    sudo powershell.exe -NoProfile -Command "Remove-MpPreference -ExclusionPath '%coreutil_rootdir%'"
) else if "%1"=="uninstall" (
    echo Uninstalling...
    sudo powershell.exe -NoProfile -Command "Add-MpPreference -ExclusionPath '%coreutil_rootdir%'"
    "%coreutil_rootdir%\installers\uninstall.exe"
    sudo powershell.exe -NoProfile -Command "Remove-MpPreference -ExclusionPath '%coreutil_rootdir%'"
) else (
    echo Uso: start.bat [first-config^|install^|uninstall]
    exit /b 1
)
