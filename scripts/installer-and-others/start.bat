@echo off
cd /d "%~dp0"

set coreutil_rootdir=%cd%
if "%1"=="install" (
    echo [INFO] Enabling Powershell Policy...
    powershell.exe -Command "Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope CurrentUser -Force"
    powershell.exe -Command "%coreutil_rootdir%\scripts\cmd\core-utils\enable-sudo-disable-installer-dectection.ps1"
    "%coreutil_rootdir%\installers\install.exe" "%coreutil_rootdir%"
) else if "%1"=="uninstall" (
    echo Uninstalling...
    "%coreutil_rootdir%\installers\uninstall.exe"
) else (
    echo Uso: start.bat [install^|uninstall]
    exit /b 1
)
