#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_packages() {
    # Git Credentials Manager
    gcm_url="https://github.com/git-ecosystem/git-credential-manager/releases/download/v__VERSION_APP__/gcm-linux-x64-__VERSION_APP__.deb"
    installer-by-url github -n git-credential-manager -u "$gcm_url" -e deb -C "sudo apt install '__DOWNLOADED_FILE_ARG__'" -o git-ecosystem -r git-credential-manager -c -R
    # Topgrade
    topgrade_url="https://github.com/topgrade-rs/topgrade/releases/download/v__VERSION_APP__/topgrade___VERSION_APP___amd64.deb"
    installer-by-url github -n topgrade -u "$topgrade_url" -e deb -C "sudo apt install '__DOWNLOADED_FILE_ARG__'" -o topgrade-rs -r topgrade -c -R
    # Powershell
    powershell_url="https://github.com/PowerShell/PowerShell/releases/download/v__VERSION_APP__/powershell___VERSION_APP__-1.deb_amd64.deb"
    installer-by-url github -n powershell -u "$powershell_url" -e deb -C "sudo apt install '__DOWNLOADED_FILE_ARG__'" -o PowerShell -r PowerShell -c -R -l
    source "$(SHELL_PROFILE_SCRIPT)"
    info-log "Install PSFzf powershell module"
    pwsh -Command "Install-Module -Name PSFzf -Scope CurrentUser -SkipPublisherCheck"
}

case "${OPERATION_ARG}" in
    packages) install_packages ;;
    *) warn-log "No GitHub Packages operation to process." ;;
esac
