#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_snap() {
    base_packages=("snapd")
    for base_package in "${base_packages[@]}"; do
        evalc "sudo apt install $base_package -y"
    done
}

case "${OPERATION_ARG}" in
    manager) install_snap ;;
    *) warn-log "No Snap operation to process." ;;
esac
