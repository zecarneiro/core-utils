#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2155
# shellcheck disable=SC2164
# shellcheck disable=SC2046
# shellcheck disable=SC2140
# shellcheck disable=SC2002
# shellcheck disable=SC1091

# ---------------------------------------------------------------------------- #
#                                      APT                                     #
# ---------------------------------------------------------------------------- #
function __install_apt {
    infolog "Update APT Repository"
    evaladvanced "sudo apt update"

    __install_apt_package 0
    __install_apt_repositories

    infolog "Update APT Repository for added repository"
    evaladvanced "sudo apt update"

    # Install Packages that depends on the repositories
    __install_apt_package 1

    infolog "Upgrade APT Package"
    evaladvanced "sudo apt upgrade -y"
}

function __install_apt_repositories {
    infolog "Add Multiverse and Universe Respositories"
    evaladvanced "sudo add-apt-repository universe -y"
    evaladvanced "sudo add-apt-repository multiverse -y"

    
}

function __install_apt_package {
    local step="$1"
    local package_list=()
    if [[ "${step}" == "1" ]]; then
        package_list=(powershell)
        addalias "powershell" "pwsh -nologo"
    else
        # Base apps
        package_list=(software-properties-common apt-transport-https wget curl inkscape git zenity libnotify-bin ubuntu-restricted-extras dos2unix fzf lsb-release)
        # Archive manager
        package_list+=(file-roller packagekit arj lzip lzma lzop ncompress rpm2cpio rzip sharutils unace unalz unar p7zip-full p7zip-rar unrar zip unzip rar uudeview mpack cabextract)
    fi
    # Install all apps
    for package_name in "${package_list[@]}"; do
        if [ $(commandexists "$package_name") == false ]; then
            infolog "Install $package_name"
            evaladvanced "sudo apt install $package_name -y"
        fi
    done

    # Install installers
    if [[ "${step}" == "0" ]]; then
        if [ $(commandexists "gdebi") == false ]; then
            infolog "Install Gdebi(DEB installers)"
            evaladvanced "sudo apt install gdebi -y"
        fi
        if [ $(commandexists "alien") == false ]; then
            infolog "Install Alien(RPM Installer)"
            evaladvanced "sudo apt install alien -y"
        fi
    fi
}

# ---------------------------------------------------------------------------- #
#                                    FLATPAK                                   #
# ---------------------------------------------------------------------------- #

# ---------------------------------------------------------------------------- #
#                                     SNAP                                     #
# ---------------------------------------------------------------------------- #


# ---------------------------------------------------------------------------- #
#                                   APPIMAGE                                   #
# ---------------------------------------------------------------------------- #


# ---------------------------------------------------------------------------- #
#                                    OTHERS                                    #
# ---------------------------------------------------------------------------- #
function __install_deb_get_package {
    evaladvanced "sudo deb-get install gcm" # git-credential-manager
}
