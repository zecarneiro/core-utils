#!/usr/bin/env bash

OPERATION_ARG="$1"; shift
FORCE_ARG="$1"; shift
SCRIPTS_APPS_DIR_ARG="$*"
SCRIPT_FORCE_ARG=""
if [[ "${FORCE_ARG}" == "true" ]]; then
    SCRIPT_FORCE_ARG="-f"
fi

function install_scripts_apps() {
    for scriptApp in "$SCRIPTS_APPS_DIR_ARG"/*; do
        if [ -f "$scriptApp" ]; then
            script-manager-cu install "$scriptApp" $SCRIPT_FORCE_ARG
        fi
    done
}

function install_alias() {
    local -A alias_map
    alias_map=(
        ["restart-pipewire"]="systemctl --user restart pipewire.service"
        ["zenity"]="$(whichc zenity -n) 2>/dev/null"
        ["now"]="date"
        ["powershell"]="pwsh"
        ["pause"]="echo -n \"Press Enter to continue...: \"; read var_name"
        ["cls"]="clear"
        ["start-menu-refresh"]="sudo update-desktop-database"
        ["gearlever"]="flatpak run it.mijorus.gearlever"
        ["gcm"]="git-credential-manager"
    )
    for alias_key in "${!alias_map[@]}"; do
        alias_value="${alias_map[$alias_key]}"
        alias-manager-cu -n "$alias_key" -c "$alias_value" $SCRIPT_FORCE_ARG
    done
}

function install_menus_entries() {
    create-menu-entry -n "Powershell" -e "$(whichc pwsh -n) -nologo" -t
    create-menu-entry -n "ChangeDNS" -e "$(whichc bash -n) -i -c '$(whichc change-dns -n)'" -t
    create-menu-entry -n "Update System CU" -e "bash -i -c '$(which system-upgrade); $(which pause)'" -i "utilities-terminal" -c "ConsoleOnly;System;" -t
    create-menu-entry -n "Cleanup System CU" -e "bash -i -c '$(which system-cleanup); $(which pause)'" -i "utilities-terminal" -c "ConsoleOnly;System;" -t
}

case "${OPERATION_ARG}" in
    scripts-apps) install_scripts_apps ;;
    alias) install_alias ;;
    menu-entries) install_menus_entries ;;
    *) warn-log "No Scripts apps or alias operation to process." ;;
esac
