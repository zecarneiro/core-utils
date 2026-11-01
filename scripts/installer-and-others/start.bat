@echo off
cd /d "%~dp0"

set coreutil_rootdir=%cd%
if "%1"=="first-config" (
    echo [INFO] Enabling Powershell Policy...
    powershell.exe -Command "Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope CurrentUser -Force"
    powershell.exe -Command "%coreutil_rootdir%\scripts\cmd\core-utils\enable-sudo.ps1"
    powershell.exe -Command "%coreutil_rootdir%\scripts\apps\pwsh\windows-policies-processor.ps1"
    sudo powershell.exe -Command "Restart-Computer -Force"
) else if "%1"=="install" (
    "%coreutil_rootdir%\installers\install.exe" "%coreutil_rootdir%"
) else if "%1"=="uninstall" (
    echo Uninstalling...
    "%coreutil_rootdir%\installers\uninstall.exe"
) else (
    echo Uso: start.bat [first-config^|install^|uninstall]
    exit /b 1
)
