#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_appimage() {
    base_packages=("libfuse2" "libfuse2t64")
    for base_package in "${base_packages[@]}"; do
        evalc "sudo apt install $base_package -y"
    done
}

function install_packages() {
    mdview_url="https://github.com/c3er/mdview/releases/download/v__VERSION_APP__/mdview-__VERSION_APP__-x86_64.AppImage"
    installer-by-url github -n mdview -u "$mdview_url" -e AppImage -C "appimage-manager-cu --install '__DOWNLOADED_FILE_ARG__'" -o c3er -r mdview -c
}

case "${OPERATION_ARG}" in
    manager) install_appimage ;;
    packages) install_packages ;;
    *) warn-log "No AppImage operation to process." ;;
esac
