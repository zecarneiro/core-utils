#!/usr/bin/env bash

OPERATION_ARG="$1"
OS_VERSION_ARG="$2"

function update_apt() {
    evalc "sudo apt update"
}

function install_apt() {
    base_packages=(
        "apt-transport-https"
		"wget"
		"curl"
		"libnotify-bin"
		"lsb-release"
		"openssl"
		"software-properties-common"
		"software-properties-gtk"
    )
    separator-log
    update_apt
    evalc "sudo add-apt-repository universe -y"
    evalc "sudo add-apt-repository multiverse -y"
    update_apt
    for base_package in "${base_packages[@]}"; do
        evalc "sudo apt install $base_package -y"
    done
}

function install_packages() {
    base_packages=(
        "zenity"
		"dos2unix"
		"git"
		"libsecret-1-0"         # For Git GCM
		"libsecret-1-dev"    # For Git GCM
		"file-roller"
		"packagekit"
		"arj"
		"lzip"
		"lzma"
		"lzop"
		"ncompress"
		"rpm2cpio"
		"rzip"
		"sharutils"
		"unace"
		"unalz"
		"unar"
		"p7zip-full"
		"7zip"
		"7zip-rar"
		"unrar"
		"zip"
		"unzip"
		"rar"
		"uudeview"
		"mpack"
		"cabextract"
		"gnome-disk-utility"
    )
    separator-log
    for base_package in "${base_packages[@]}"; do
        evalc "sudo apt install $base_package -y"
    done
}

case "${OPERATION_ARG}" in
    manager) install_apt ;;
    packages) install_packages ;;
    *) warn-log "No APT operation to process." ;;
esac
