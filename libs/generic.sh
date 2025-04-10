#!/usr/bin/env bash

function __show_install_message_question {
    read -p "Do you want to install ${1}? (y/N): " userInput
    echo "$userInput"
}

function __create_dirs {
    local dirs=("$CONFIG_DIR" "$OTHER_APPS_DIR" "$USER_BIN_DIR" "$USER_STARTUP_DIR" "$USER_TEMP_DIR" "$TEMP_DIR")
    for dir in "${dirs[@]}"; do
        if [[ ! -d "${dir}" ]]; then
            mkdir "$dir"
            echo "Created directory: $dir"
        fi
    done
}

function __exit_script {
    exit 0
}

function create-profile-file-bash {
    if [ ! -f "$MY_SHELL_PROFILE" ]; then
        touch "$MY_SHELL_PROFILE"
    fi
    if [ "$(filecontain "$MY_SHELL_PROFILE" "$MY_CUSTOM_SHELL_PROFILE")" == false ]; then
        writefile "$MY_SHELL_PROFILE" "source '$MY_CUSTOM_SHELL_PROFILE'" -append
    fi
    if [ ! -f "$MY_CUSTOM_SHELL_PROFILE" ]; then
        infolog "Creating Bash Script profile to run when bash start: $MY_CUSTOM_SHELL_PROFILE"
        touch "$MY_CUSTOM_SHELL_PROFILE"
    fi
}

function install-profile-scripts {
    local shellScriptsInstallDir="${OTHER_APPS_DIR}/shell-scripts"
    local scriptsList=($(find "$shellScriptsInstallDir" -maxdepth 1 -type f -name '*.sh'))
    if [[ "${onlyProfile}" == false ]]; then
        infolog "Install core-utils scripts release package"
        rm -rf "$shellScriptsInstallDir"
        cpdir "$SHELL_SCRIPT_DIR" "$shellScriptsInstallDir"
    fi
    # Add bash profiles
    for script in "${scriptsList[@]}"; do
        local data="source '$script'"
        if [ "$(filecontain "$MY_CUSTOM_SHELL_PROFILE" "$data")" == false ]; then
            writefile "$MY_CUSTOM_SHELL_PROFILE" "$data" -append
        fi
    done
}
