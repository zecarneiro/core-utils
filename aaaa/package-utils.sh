#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2155
# shellcheck disable=SC2164
# shellcheck disable=SC2046
# shellcheck disable=SC2140
# shellcheck disable=SC2002

# ---------------------------------------------------------------------------- #
#                                   FUNCTIONS                                  #
# ---------------------------------------------------------------------------- #
function install_all_packages() {
    install_apt
    install_flatpak
    install_snap
    install_deb_get
    install_pacstall
    install_appimage
    install_installers
}

function install_all_base_packages() {
    install_base_apt_package
    install_base_appimage_package
}

# ---------------------------------------------------------------------------- #
#                                 PACKAGES AREA                                #
# ---------------------------------------------------------------------------- #
function install_apt {
    log "\nAdd Multiverse and Universe Respositories"
    evaladvanced "sudo add-apt-repository universe -y"
    evaladvanced "sudo add-apt-repository multiverse -y"

    log "\nUpdate APT Repository"
    evaladvanced "sudo apt update"

    log "\nUpgrade APT Package"
    evaladvanced "sudo apt upgrade -y"
}

function install_flatpak {
    if [ $(commandexists "flatpak") == false ]; then
        log "\nInstall Flatpak"
        evaladvanced "sudo apt install gnome-software gnome-software-plugin-flatpak xdg-desktop-portal-gtk flatpak -y"
        evaladvanced "flatpak remote-add --if-not-exists flathub 'https://flathub.org/repo/flathub.flatpakrepo'"
    fi
}

function install_snap {
    if [ $(commandexists "snap") == false ]; then
        log "\nInstall Snap"
        evaladvanced "sudo apt install snapd-xdg-open snapd -y"
    fi
}

function install_pacstall {
    if [ $(commandexists "pacstall") == false ]; then
        log "\nInstall Pacstall"
        download --url "https://pacstall.dev/q/install" --file "$APPS_BIN_DIR/install_pacstall.sh"
        evaladvanced "chmod +x \"$APPS_BIN_DIR/install_pacstall.sh\""
        evaladvanced "sudo \"$APPS_BIN_DIR/install_pacstall.sh\""
    fi
}

function install_appimage {
    log "\nEnable AppImage Support in Ubuntu"
    evaladvanced "sudo apt install libfuse2 -y"
    evaladvanced "sudo apt install libfuse2t64 -y"
}

function install_installers {
    if [ $(commandexists "gdebi") == false ]; then
        log "\nInstall Gdebi(DEB installers)"
        evaladvanced "sudo apt install gdebi -y"
    fi
    if [ $(commandexists "alien") == false ]; then
        log "\nInstall Alien(RPM Installer)"
        evaladvanced "sudo apt install alien -y"
    fi
}

# ---------------------------------------------------------------------------- #
#                              BASE PACKAGES AREA                              #
# ---------------------------------------------------------------------------- #
function install_base_apt_package {
    local package_list=(software-properties-common apt-transport-https wget curl inkscape git zenity libnotify-bin)
    for package_name in "${package_list[@]}"; do
        if [ $(commandexists "$package_name") == false ]; then
            log "\nInstall $package_name"
            evaladvanced "sudo apt install $package_name -y"
        fi
    done
}

function install_base_appimage_package {
    infolog "Base Appimage package is empty"
}
