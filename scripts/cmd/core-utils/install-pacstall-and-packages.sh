#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_pacstall() {
    sudo bash -c "$(curl -fsSL https://pacstall.dev/q/install)" <<< $'n\n'
}

function install_packages() {
    base_packages=("fzf-bin")
    for base_package in "${base_packages[@]}"; do
        evalc "sudo pacstall -P -I $base_package"
    done
}

case "${OPERATION_ARG,,}" in
    manager) install_pacstall ;;
    packages) install_packages ;;
    *) warn-log "No Pacstall operation to process." ;;
esac
