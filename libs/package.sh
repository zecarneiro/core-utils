#!/usr/bin/env bash
# shellcheck disable=SC1090

function __set_user_bin_dir() {
    local user_bin_dir="$HOME/.local/bin"
    __print "Set User bin dir"
    __eval "mkdir -p '$user_bin_dir'"
}

function __install_scripts() {
    __eval "sudo cp '$__SCRIPTS_DIR__/pipc.sh' '/usr/local/bin/pipc'"
    __eval "sudo chmod +x /usr/local/bin/pipc"

    __eval "sudo cp '$__SCRIPTS_DIR__/appimage-manager.sh' '/usr/local/bin/appimage-manager'"
    __eval "sudo chmod +x /usr/local/bin/appimage-manager"

    __eval "sudo cp '$__SCRIPTS_DIR__/sudoexec.sh' '/usr/local/bin/sudoexec'"
    __eval "sudo chmod +x /usr/local/bin/sudoexe"
}

function __install_apt_and_packages() {
    __print "Install APT and Packages"
    __eval "sudo apt update"
    __eval "sudo add-apt-repository universe -y"
    __eval "sudo add-apt-repository multiverse -y"

    __eval ". '$__SHELL_FILE__'" true
    . "$__SHELL_FILE__"

    __eval "sudo apt install apt-transport-https -y"
    __eval "sudo apt install wget -y"
    __eval "sudo apt install curl -y"
    __eval "sudo apt install libnotify-bin -y"
    __eval "sudo apt install lsb-release -y"
    __eval "sudo apt install fzf -y"
    __eval "sudo apt install zenity -y"
    __eval "sudo apt install dos2unix -y"

    __eval ". '$__SHELL_FILE__'" true
    . "$__SHELL_FILE__"

    __print "Install microsoft repository"
    local microsoft_repo_deb_file="/tmp/packages-microsoft-prod.deb"
    # Get the version of Ubuntu
    # shellcheck disable=SC1091
    source /etc/os-release
    # Download the Microsoft repository keys
    wget -O "${microsoft_repo_deb_file}" "https://packages.microsoft.com/config/ubuntu/${VERSION_ID}/packages-microsoft-prod.deb" -q --show-progress
    # Register the Microsoft repository keys
    __eval "sudo dpkg -i '$microsoft_repo_deb_file'"
    # Delete the Microsoft repository keys file
    __eval "rm '$microsoft_repo_deb_file'"

    __eval ". '$__SHELL_FILE__'" true
    . "$__SHELL_FILE__"

    __print "Install Python"
    __eval "sudo apt install python3 -y"
    __eval "sudo apt install python-is-python3 -y"
    __eval "sudo apt install python3-pip -y"
    __eval "sudo apt install python3-venv -y"
    __eval "sudo apt install pipx -y"
    __eval "pipx ensurepath"
    __eval "pipc install pyright"
    __eval "pipc install ruff"
}


function __install_flatpak_and_packages {
    __print "Install Flatpak"
    __eval "sudo apt install gnome-software gnome-software-plugin-flatpak xdg-desktop-portal-gtk flatpak -y"
    __eval "flatpak remote-add --if-not-exists flathub 'https://flathub.org/repo/flathub.flatpakrepo'"

    __eval ". '$__SHELL_FILE__'" true
    . "$__SHELL_FILE__"

    __eval "flatpak install flathub it.mijorus.gearlever -y"
    __eval "sudo flatpak override it.mijorus.gearlever --filesystem=host"
}

function __install_snap_and_packages {
    __print "Install Snap"
    __eval "sudo apt install snapd-xdg-open snapd -y"
}

function __install_appimage_and_packages {
    __print "Enable AppImage Support in Ubuntu"
    __eval "sudo apt install libfuse2 -y"
    __eval "sudo apt install libfuse2t64 -y"

    __eval ". '$__SHELL_FILE__'" true
    . "$__SHELL_FILE__"

    __print "Install mdview"
    wget -O "/tmp/mdview.AppImage" "https://github.com/c3er/mdview/releases/download/v3.2.0/mdview-3.2.0-x86_64.AppImage" -q --show-progress
    appimage-manager --install "/tmp/mdview.AppImage"
}

function __install_pacstall_and_packages {
    __print "Install Pacstall"
    __eval "sudo bash -c \"\$(curl -fsSL https://pacstall.dev/q/install)\" <<< $'n\n'"
}

function __install_deb_get_and_packages {
    __print "Install Deb-Get"
    __eval "curl -sL https://raw.githubusercontent.com/wimpysworld/deb-get/main/deb-get | sudo -E bash -s install deb-get"

    __eval ". '$__SHELL_FILE__'" true
    . "$__SHELL_FILE__"
    
    __eval "sudo deb-get install topgrade"
}
