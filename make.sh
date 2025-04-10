#!/usr/bin/env bash
declare start=false
declare onlyProfile=false

declare SCRIPT_UTILS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
declare VERSION="$(cat "$SCRIPT_UTILS_DIR/version")"
declare SHELL_SCRIPT_DIR="${SCRIPT_UTILS_DIR}/scripts"
declare LIBS_DIR="${SCRIPT_UTILS_DIR}/libs"
declare BIN_DIR="${SCRIPT_UTILS_DIR}/bin"

if [[ "${1}" == "-s" ]]||[[ "${1}" == "--start" ]]; then
    start=true
fi
if [[ "${2}" == "--only-profile-shell" ]]; then
    onlyProfile=true
fi

# ---------------------------------------------------------------------------- #
#                                    IMPORTS                                   #
# ---------------------------------------------------------------------------- #
for script in "${SHELL_SCRIPT_DIR}"/*.sh; do
    source "$script"
done
for script in "${LIBS_DIR}"/*.sh; do
    source "$script"
done

# ---------------------------------------------------------------------------- #
#                                     MAIN                                     #
# ---------------------------------------------------------------------------- #
function usage {
    echo "Usage: make.ps1 [OPTIONS]... [STEP-VALUE]"
    echo -e "OPTIONS:
-s|--start\tProcess install and config by user"
}

function printMenu {
    echo "1. Will
    - Add important APT repository
    - Install APT packages
2. Will
    - Install FLATPAK
    - Install SNAP
    - Install PACSTALL
    - Install APPIMAGE
    - Install GDEBI, ALIEN and ZAP
3. Will
    - Install Snap and Flatpak packages
    - Install PIP/PIPX
    - Install NPM
4. Will
    - Create user powershell profile file
    - Install scripts profile
5. Will
    - Install Appimage packages
    - Install Development packages. User will decide wich to install
---
6. Exit"
}

function initProcess {
    local message="Please, restart your terminal."
    __create_dirs
    while true; do
        local option=-1
        if [[ "${onlyProfile}" == false ]]; then
            printMenu
            read -p "Insert an option: " option
        else
            option=4
        fi
        case "${option}" in
            1)
                install-apt
                warnlog "$message"
                __exit_script
            ;;
            2)
                install-flatpak
                install-snap
                install-pacstall
                install-appimage
                install-linstallers
                warnlog "$message"
                __exit_script
            ;;
            3)
                install-snap-package
                install-flatpak-package
                install-pip-pipx
                install-npm
                warnlog "$message"
                __exit_script
            ;;
            4)
                create-profile-file-bash
                install-profile-scripts
                if [[ "${onlyProfile}" == false ]]; then
                    cd "$SCRIPT_UTILS_DIR" || __exit_script
                    pwsh -command ".\make.ps1 --start --only-profile-shell"
                    warnlog "$message"
                fi
                __exit_script
            ;;
            5)
                install-appimage-packages
                install-development-package
                warnlog "$message"
                __exit_script
            ;;
            6) __exit_script ;;
            *) warnlog "Please, insert a valid option!" ;;
        esac
    done
}

function main {
    if [[ "${start}" == true ]]; then
        initProcess
    else
        usage
    fi
}
main
