# Author: José M. C. Noronha
# This script need coreutils package

# Validate core utils bin directory
$coreutilsBinDir = (Resolve-Path "$PSScriptRoot\..\..\..\bin").Path
if (!(Test-Path "$coreutilsBinDir") -or !(Test-Path "$coreutilsBinDir\appx-package-list.exe" -PathType Leaf)) {
    Write-Error "Invalid CoreUtils bin directory"
    exit 1
}
$env:Path += ";${coreutilsBinDir}"
title-log "Install Winget-CLI"
if ($PSVersionTable.PSVersion.Major -eq 7) {
    warn-log "This command does not work in PowerShell 7. You must install in Windows PowerShell."
    exit 1
}

header-log "Check for VCLibs requirement"
$requirement = (appx-package-list -f "VCLibs.140.00.UWPDesktop")
if (-Not $requirement) {
    info-log "Installing Desktop App Installer requirement"
    Try {
        $output = "$(TEMP_DIR)\VCLibs.140.00.UWPDesktop.appx"
        download -u "https://aka.ms/Microsoft.VCLibs.x64.14.00.Desktop.appx" -o "$output"
        run-bin-processor "$output"
        evalc "rmfile `"$output`""
        ok-log "VCLibs installed"
    } Catch {
        error-log "On install VCLibs.140.00.UWPDesktop"
    }
}

header-log "Check for WindowsAppRuntime requirement"
$requirement = (appx-package-list -f "Microsoft.WindowsAppRuntime.1.8")
if (-not $requirement) {
    info-log "Installing Windows App Runtime requirement"
    Try {
        $output = "$(TEMP_DIR)\WindowsAppRuntime_x64.exe"
        download -u "https://aka.ms/windowsappsdk/1.8/latest/windowsappruntimeinstall-x64.exe" -o "$output"
        run-bin-processor "$output"
        evalc "rmfile `"$output`""
        ok-log "WindowsAppRuntime installed"
    } Catch {
        error-log "On install Microsoft.WindowsAppRuntime.1.8"
    }
}

header-log "Install/Update Winget"
if (!(Get-Command winget -ErrorAction SilentlyContinue)) {
    info-log "Winget not found. Installing App Installer..."
    Try {
        $output = "$(TEMP_DIR)\AppInstaller.msixbundle"
        download -u "https://aka.ms/getwinget" -o "$output"
        run-bin-processor "$output"
        evalc "rmfile `"$output`""
        ok-log "Winget installed"
    } Catch {
        error-log "On install App Installer"
    }
} else {
    info-log "Winget already installed. Updating..."
    Try {
        evalc "sudo winget source reset --force"
        evalc "sudo winget source update"
        winget-install "Microsoft.AppInstaller"
        ok-log "Winget updated"
    } Catch {
        error-log "On update Winget"
    }
}
