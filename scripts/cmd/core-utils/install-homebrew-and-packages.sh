#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_homebrew() {
    shell_script="$(SHELL_PROFILE_SCRIPT)"
    evalc "/bin/bash -c '$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)'"
    if [ "$(cat "$shell_script" | grep -c "/home/linuxbrew")" -le 0 ]; then
        echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> "$shell_script"
        eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
    fi
}

function install_packages() {
    info-log "Nothing to process..."
}



case "${OPERATION_ARG}" in
    manager) install_homebrew ;;
    packages) install_packages ;;
    *) warn-log "No Deb-Get operation to process." ;;
esac
