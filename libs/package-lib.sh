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

    infolog "Install microsoft repository"
    local repoDebFile="$TEMP_DIR/packages-microsoft-prod.deb"
    source /etc/os-release # Get the version of Ubuntu
    download --url "https://packages.microsoft.com/config/ubuntu/${VERSION_ID}/packages-microsoft-prod.deb" --file "$repoDebFile" # Download the Microsoft repository keys
    evaladvanced "sudo dpkg -i '$repoDebFile'" # Register the Microsoft repository keys
    evaladvanced "rm '$repoDebFile'" # Delete the Microsoft repository keys file
}

function __install_apt_package {
    local step="$1"
    local package_list=()
    if [[ "${step}" == "1" ]]; then
        package_list=(powershell)
    else
        # Base apps
        package_list=(software-properties-common apt-transport-https wget curl inkscape git zenity libnotify-bin ubuntu-restricted-extras dos2unix fzf)
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
function __install_flatpak {
    if [ $(commandexists "flatpak") == false ]; then
        infolog "Install Flatpak"
        evaladvanced "sudo apt install gnome-software gnome-software-plugin-flatpak xdg-desktop-portal-gtk flatpak -y"
        evaladvanced "flatpak remote-add --if-not-exists flathub 'https://flathub.org/repo/flathub.flatpakrepo'"
    fi
}
function __install_flatpak_package {
    evaladvanced "flatpak install flathub it.mijorus.gearlever -y"
    evaladvanced "sudo flatpak override it.mijorus.gearlever --filesystem=host"
    addalias "gearlever" "flatpak run it.mijorus.gearlever"
}

# ---------------------------------------------------------------------------- #
#                                     SNAP                                     #
# ---------------------------------------------------------------------------- #
function __install_snap {
    if [ $(commandexists "snap") == false ]; then
        infolog "Install Snap"
        evaladvanced "sudo apt install snapd-xdg-open snapd -y"
    fi
}

# ---------------------------------------------------------------------------- #
#                                   APPIMAGE                                   #
# ---------------------------------------------------------------------------- #
function __install_appimage {
    infolog "Enable AppImage Support in Ubuntu"
    evaladvanced "sudo apt install libfuse2 -y"
    evaladvanced "sudo apt install libfuse2t64 -y"
}

function __install_appimage_packages {
    appimage-install --url "https://github.com/c3er/mdview/releases/download/v3.2.0/mdview-3.2.0-x86_64.AppImage"
}

# ---------------------------------------------------------------------------- #
#                                    OTHERS                                    #
# ---------------------------------------------------------------------------- #
function __install_pacstall {
    if [ $(commandexists "pacstall") == false ]; then
        infolog "Install Pacstall"
        sudo bash -c "$(curl -fsSL https://pacstall.dev/q/install)"
    fi
}
