#!/usr/bin/env bash
# shellcheck disable=SC2001

declare user_bin_dir="$HOME/.local/bin"
declare user_desktop_files_dir="$HOME/.local/share/applications"
declare line_tag="# Inserted by coreutils, zecarneiro"

if ! command -v flatpak &>/dev/null; then
    echo "[ERROR] Please, install flatpak!"
    exit 1
fi

if (( $(flatpak list | grep -c it.mijorus.gearlever) <= 0 )); then
    echo "[ERROR] Please, install Gearlever from flatpak"
    exit 1
fi

__install() {
    local file="$1"
    if [[ ! -f "$file" ]]; then
        echo "[ERROR] Invalid given file: ${file}"
        exit 1
    fi
    echo "### INSTALL '${file}' ###"
    cp "$file" "${file}-copy"
    flatpak run it.mijorus.gearlever --integrate "${file}-copy" -y
    sleep 1
    __update_commands
}

__update_commands() {
    local app_installed_list=($(flatpak run it.mijorus.gearlever --list-installed | awk '{print $NF}'))
    __uninstall_commands
    echo "### INSTALL ALL COMMANDS FOR EACH APPS ###"
    for app in "${app_installed_list[@]}"; do
        error_msg="[ERROR] Failed on install command for app: ${app}"
        nosandbox_arg="--no-sandbox"
        app_command=$(basename "$app")
        app_command=$(echo "$app_command" | sed "s/.appimage//g")
        app_command_file="${user_bin_dir}/${app_command}"
        app_desktop_file="${user_desktop_files_dir}/${app_command}.desktop"
        if [ ! -f "${app_desktop_file}" ]; then
            echo "${error_msg}"
            continue
        fi
        if [[ $(cat "$app_desktop_file" | grep -c "no-sandbox") -eq 0 ]]; then
            nosandbox_arg=""
        fi
        echo -e "#!/bin/bash\n${line_tag}\n${app} ${nosandbox_arg} \$*" | tee "$app_command_file" >/dev/null
        if [ -f "$app_command_file" ]; then
            chmod +x "$app_command_file"
            echo "[OK] Installed command: $app_command"
        else
            echo "${error_msg}"
        fi
    done
}

__uninstall_commands() {
    echo "### UNINSTALL ALL INSTALLED COMMANDS ###"
    find "$user_bin_dir" -type f -exec grep -q -l "${line_tag}" '{}' \; -delete -print
}

main() {
    case "${1}" in
        -i|--install) __install "$2" ;;
        -c|--update-commands) __update_commands ;;
        -u|--uninstall-commands) __uninstall_commands ;;
        *)
            echo "$0 [-i|--install] FILE | [-c|--update-commands] | [-u|--uninstall-commands]"
        ;;
    esac
    
}
main "$@"
