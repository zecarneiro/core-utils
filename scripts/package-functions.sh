#!/usr/bin/env bash
# Author: José M. C. Noronha

# ---------------------------------------------------------------------------- #
#                                      NPM                                     #
# ---------------------------------------------------------------------------- #
function npmupgrade {
    evaladvanced "npm outdated -g"
    evaladvanced "npm update -g"
}
function npmlist {
    local filter="$1"
    local command_to_run="npm list -g --depth=0"
    if [[ -n "${filter}" ]]; then
        command_to_run="${command_to_run} | grep ${filter}"
    fi
    evaladvanced "${command_to_run}"
}

# ---------------------------------------------------------------------------- #
#                                      APT                                     #
# ---------------------------------------------------------------------------- #
function aptupgrade {
    evaladvanced "sudo apt update"
    evaladvanced "sudo apt upgrade -y"
}
function aptuninstall {
    evaladvanced "sudo apt purge --autoremove '$1' -y"
    evaladvanced "aptclean"
}
function aptclean {
    evaladvanced 'sudo apt autoremove -y'
    evaladvanced 'sudo apt autopurge -y'
    evaladvanced 'sudo apt autoclean -y'
    evaladvanced 'sudo apt clean --dry-run'
}
function aptlist {
    local filter="$1"
    local command_to_run="apt list --installed"
    if [[ -n "${filter}" ]]; then
        command_to_run="${command_to_run} | grep ${filter}"
    fi
    evaladvanced "${command_to_run}"
}

# ---------------------------------------------------------------------------- #
#                                    FLATPAK                                   #
# ---------------------------------------------------------------------------- #
function flatpakupgrade {
    evaladvanced "flatpak update -y"
}
function flatpakuninstall {
    evaladvanced "flatpak uninstall --delete-data '$1' -y"
}
function flatpakclean {
    evaladvanced 'flatpak uninstall --unused -y'
    evaladvanced 'sudo rm -rfv /var/tmp/flatpak-cache*'
}
function flatpaklist {
    local filter="$1"
    local command_to_run="flatpak list"
    if [[ -n "${filter}" ]]; then
        command_to_run="${command_to_run} | grep ${filter}"
    fi
    evaladvanced "${command_to_run}"
}

# ---------------------------------------------------------------------------- #
#                                     SNAP                                     #
# ---------------------------------------------------------------------------- #
function snapupgrade {
    evaladvanced 'sudo snap refresh'
}
function snapuninstall {
    local configDir="$HOME/snap/$1"
    local configSystemDir="/snap/$1"
    evaladvanced "sudo snap remove --purge '$1'"
    evaladvanced "snap saved"
    read -p "Insert the number on the line of App(ENTER TO SKIP): " userInput
    if [[ -n "${userInput}" ]]; then
        evaladvanced "sudo snap forget ${userInput}"
    fi
    if [ -d "$configDir" ]; then
        evaladvanced "rm -rf \"$configDir\""
    fi
    if [ -d "$configSystemDir" ]; then
        evaladvanced "sudo rm -rf \"$configSystemDir\""
    fi
}
function snapclean {
    evaladvanced "sudo sh -c 'rm -rf /var/lib/snapd/cache/*'"
    warnlog "Removes old revisions of snaps"
    read -n 1 -s -r -p "Please, CLOSE ALL SNAPS BEFORE RUNNING THIS. PRESS ANY KEY TO CONTINUE"
    log
    LANG=en_US.UTF-8 snap list --all | awk '/disabled/{print $1, $3}' | while read snapname revision; do
        evaladvanced "sudo snap remove \"$snapname\" --revision=\"$revision\""
    done
}
function snaplist {
    local filter="$1"
    local command_to_run="snap list"
    if [[ -n "${filter}" ]]; then
        command_to_run="${command_to_run} | grep ${filter}"
    fi
    evaladvanced "${command_to_run}"
}

# ---------------------------------------------------------------------------- #
#                               UPDATERS SCRIPTS                               #
# ---------------------------------------------------------------------------- #
function installupdater {
    # shellcheck disable=SC2116
    # shellcheck disable=SC2155
    local updater_dir="$(echo "$HOME")/.otherapps/updaters"
    local updater_script="$1"
    # shellcheck disable=SC2155
    local scriptname=$(basename "$updater_script")
    mkdir -p "$updater_dir"
    infolog "Installing '$scriptname' '$updater_script'"
    cp "$updater_script" "$updater_dir"
    chmod -R 777 "$updater_dir"
    oklog "Done"
}
function updatersupgrade {
    local scriptname="$1"
    local currentdir="$PWD"
    # shellcheck disable=SC2116
    # shellcheck disable=SC2155
    local updater_dir="$(echo "$HOME")/.otherapps/updaters"
    if [ -d "$updater_dir" ]; then
        for script in "$updater_dir"/*; do
            # shellcheck disable=SC2155
            local updatername=$(basename "$script")
            if [[ -z "${scriptname}" ]]||[[ "${scriptname}" == "${updatername}" ]]; then
                promptlog "$script"
                # shellcheck disable=SC1090
                . "$script"
            fi
        done
    fi
    # shellcheck disable=SC2164
    cd "$currentdir"
}

# ---------------------------------------------------------------------------- #
#                                SYSTEM PACKAGES                               #
# ---------------------------------------------------------------------------- #
alias systemupgrade="npmupgrade; log; aptupgrade; log; flatpakupgrade; log; snapupgrade; log; debgetupgrade; updatersupgrade"
alias systemclean="aptclean; log; flatpakclean; log; snapclean; log; debgetclean"
