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
#                                   PACSTALL                                   #
# ---------------------------------------------------------------------------- #
function pacstalllist {
    evaladvanced "pacstall --list"
}
function pacstallupgrade {
    evaladvanced 'sudo pacstall -Up'
}

# ---------------------------------------------------------------------------- #
#                               UPDATERS SCRIPTS                               #
# ---------------------------------------------------------------------------- #
function installupdater {
    # shellcheck disable=SC2116
    # shellcheck disable=SC2155
    local updater_dir="$OTHER_APPS_DIR/updaters"
    local updater_script="$1"
    # shellcheck disable=SC2155
    local scriptname=$(basename "$updater_script")
    mkdir -p "$updater_dir"
    infolog "Installing '$scriptname' '$updater_script'"
    cp "$updater_script" "$updater_dir"
    chmod -R 777 "$updater_dir"
    updatersupgrade "$scriptname"
    oklog "Done"
}
function updatersupgrade {
    local scriptname="$1"
    local currentdir="$PWD"
    # shellcheck disable=SC2116
    # shellcheck disable=SC2155
    local updater_dir="$OTHER_APPS_DIR/updaters"
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
function systemupgrade {
    npmupgrade
    log; aptupgrade
    log; flatpakupgrade
    log; snapupgrade
    log; pacstallupgrade
    log; evaladvanced "pipx upgrade-all"
    updatersupgrade ""
}
function systemclean {
    aptclean
    log; flatpakclean
    log; snapclean
    log; debgetclean
}

# ---------------------------------------------------------------------------- #
#                                   APPIMAGE                                   #
# ---------------------------------------------------------------------------- #
appimage-install() {
    local url=""
    local file="$(mktemp -ut appimagetoinstallXXX).appimage"
    while [ "$#" -ne 0 ] ; do
        case "${1}" in
            --url) url="$2"; shift 2 ;;
            --file)
                file="$2"; shift 2
            ;;
			--help) log "Will install appimage with gearlever\n>>> appimage-install --url URL|--file FILE"; return ;;
            *) shift ;;
        esac
    done
    if [[ -n "$url" ]]; then
        file="$(basename "$file")"
        file="$USER_TEMP_DIR/$file"
        download --url "$url" --file "$file"
    else
        if [[ ! -f "$file" ]]; then
            errorlog "Invalid given file"
            return
        fi
    fi
    flatpak run it.mijorus.gearlever --integrate "$file" -y
    headerlog "Update commands with news Gearlever apps"
    appimage-command-manager --
    if [[ -f "$file" ]]; then
        rm "$file"
    fi
}

appimage-command-manager() {
    local installedAppsDir=""
    local uninstallOperation=false
    local listInstalled=false
    local binDir="$HOME/.local/bin"
    local userDesktopFilesDir="$HOME/.local/share/applications"
    local lineTag="# Inserted by core-utils, zecarneiro"
    local listApps=($(flatpak run it.mijorus.gearlever --list-installed))
    if [ ! -d "$binDir" ]; then
        errorlog "Not found: '$binDir'"
        exit 1
    fi
    while [ "$#" -ne 0 ] ; do
        case "${1}" in
            --installed-dir) installedAppsDir="$2"; shift 2 ;;
            --uninstall) uninstallOperation=true; shift ;;
            --list-installed) listInstalled=true; shift ;;
			--help) log "gearlever-bin-manager |--installed-dir 'GEARLEVER_INSTALLED_APPS_DIR' |--uninstall |--list-installed"; return ;;
            *) shift ;;
        esac
    done
    if [[ "${listInstalled}" == true ]]; then
        find "$binDir" -type f -print0 | xargs -0 grep -l "$lineTag"
        return
    fi

    # Uninstall all commands
    find "$binDir" -type f -exec grep -q -l "$lineTag" '{}' \; -delete
    if [[ "${uninstallOperation}" == true ]]; then
        oklog "Done."
        return
    fi
    
    # Install all commands
    if [[ -z "${installedAppsDir}" ]]; then
        installedAppsDir="$HOME/AppImages"
    fi
    for app in "${listApps[@]}"; do
        local commandName="$(filename "$(basename "$app")")"
        if [[ "${app}" == *"$installedAppsDir"* ]]; then
            local commandFile="$binDir/$commandName"
            local desktopFile="$(find "$userDesktopFilesDir" -type f -print0 | xargs -0 grep -l "${app}")"
            local nosandboxArg="--no-sandbox"
            if [[ $(cat "$desktopFile" | grep -c "no-sandbox") -eq 0 ]]; then
                nosandboxArg=""
            fi
            echo -e "#!/bin/bash\n${lineTag}\n${app} ${nosandboxArg} \$*" | tee "$commandFile" >/dev/null
            chmod +x "$commandFile"
            oklog "Installed command: $commandName"
        fi
    done
    oklog "Done."
}

# ---------------------------------------------------------------------------- #
#                                    DEB GET                                   #
# ---------------------------------------------------------------------------- #
function debgetlist {
    evaladvanced "deb-get list"
}
function debgetupgrade {
    evaladvanced 'sudo deb-get upgrade'
}
function debgetuninstall {
    evaladvanced "sudo deb-get purge $1"
}
function debgetclean {
    evaladvanced "sudo deb-get clean"
}
