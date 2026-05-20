#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_flatpak() {
    base_packages=(
        "xdg-desktop-portal-gtk"
        "flatpak"
    )
    for base_package in "${base_packages[@]}"; do
        evalc "sudo apt install $base_package -y"
    done
    evalc "flatpak remote-add --if-not-exists flathub 'https://dl.flathub.org/repo/flathub.flatpakrepo'"
    evalc "flatpak remote-add --user --if-not-exists flathub 'https://dl.flathub.org/repo/flathub.flatpakrepo'"
}

function install_packages() {
    base_packages=("it.mijorus.gearlever")
    for base_package in "${base_packages[@]}"; do
        flatpak-install "$base_package"
    done
}

case "${OPERATION_ARG}" in
    manager) install_flatpak ;;
    packages) install_packages ;;
    *) warn-log "No Flatpak operation to process." ;;
esac
