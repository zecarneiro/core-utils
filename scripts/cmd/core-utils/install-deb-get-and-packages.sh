#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_deb_get() {
    evalc "curl -sL https://raw.githubusercontent.com/wimpysworld/deb-get/main/deb-get | sudo -E bash -s install deb-get"
    source "$(SHELL_PROFILE_SCRIPT)"
    evalc "deb-get update"
}

function install_packages() {
    base_packages=("powershell" "gcm" "topgrade")
    for base_package in "${base_packages[@]}"; do
        evalc "sudo deb-get install $base_package"
    done
    echo
    source "$(SHELL_PROFILE_SCRIPT)"
    info-log "Install PSFzf powershell module"
    pwsh -Command "Install-Module -Name PSFzf -Scope CurrentUser -SkipPublisherCheck"
}



case "${OPERATION_ARG}" in
    manager) install_deb_get ;;
    packages) install_packages ;;
    *) warn-log "No Deb-Get operation to process." ;;
esac
