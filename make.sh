#!/usr/bin/env bash
declare start=false
declare onlyProfile=false

declare SCRIPT_UTILS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
declare VERSION="$(cat "$SCRIPT_UTILS_DIR/version")"
declare SHELL_SCRIPT_DIR="${SCRIPT_UTILS_DIR}/scripts"
declare LIBS_DIR="${SCRIPT_UTILS_DIR}/libs"
declare BIN_DIR="${SCRIPT_UTILS_DIR}/bin"
declare IMAGES_DIR="${SCRIPT_UTILS_DIR}/images"
declare USER_OPTION_INSERTED_FILE="${SCRIPT_UTILS_DIR}/user-option-inserted"

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
function manageUserOption {
    local option="$1"
    local isGet="$2"
    local doneRes=""
    if [[ "${isGet}" == true ]]; then
        if [[ $(filecontain "$USER_OPTION_INSERTED_FILE" "$option") == true ]]; then
            doneRes=" - Done"
        fi
    else
        writefile "$USER_OPTION_INSERTED_FILE" "$option" true
    fi
    echo "$doneRes"
}

function usage {
    echo "Usage: make.ps1 [OPTIONS]... [STEP-VALUE]"
    echo -e "OPTIONS:
-s|--start\tProcess install and config by user"
}

function printMenu {
    echo "1. Will$(manageUserOption 1 true)
    - Add important APT repository
    - Install APT packages
2. Will$(manageUserOption 2 true)
    - Install APPIMAGE
    - Install FLATPAK
    - Install SNAP
    - Install PACSTALL
3. Will$(manageUserOption 3 true)
    - Install Flatpak packages
    - Install Appimage packages
4. Will$(manageUserOption 4 true)
    - Create user powershell profile file
    - Install scripts profile
5. Will$(manageUserOption 5 true)
    - Install Development packages. User will decide wich to install
    - Start all necessary configurations
6. Will(Optional)$(manageUserOption 6 true)
    - Define/Change default system dirs. Like Documents, Images, etc
---
7. Exit"
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
                __install_apt
                warnlog "$message"
                manageUserOption 1 false
                __exit_script
            ;;
            2)
                __install_appimage
                __install_flatpak
                __install_snap
                __install_pacstall
                warnlog "$message"
                manageUserOption 2 false
                __exit_script
            ;;
            3)
                __install_flatpak_package
                __install_appimage_packages
                warnlog "$message"
                manageUserOption 3 false
                __exit_script
            ;;
            4)
                __create_profile_file_bash
                __install_profile_scripts
                if [[ "${onlyProfile}" == false ]]; then
                    cd "$SCRIPT_UTILS_DIR" || __exit_script
                    pwsh -command ".\make.ps1 --start --only-profile-shell"
                    warnlog "$message"
                fi
                manageUserOption 4 false
                __exit_script
            ;;
            5)
                __install_development_package
                __config_all
                warnlog "$message"
                manageUserOption 5 false
                __exit_script
            ;;
            6)
                __define_default_system_dir
                manageUserOption 6 false
            ;;
            7) __exit_script ;;
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
