#!/usr/bin/env bash
# Author: José M. C. Noronha
# Some code has source: https://github.com/ChrisTitusTech/mybash

shopt -s expand_aliases # Enable alias to use on bash script

function reboot {
    local message="Will be restart PC. Continue(y/N)? "
    read -p "$message" userInput
    if [[ "${userInput}" == "Y" ]] || [[ "${userInput}" == "y" ]]; then
        sudo "$(which shutdown)" -r now
    fi
}
function shutdown {
    read -p "Will be shutdown PC. Continue(Y/n)? " userInput
    if [[ -z "${userInput}" ]] || [[ "${userInput}" == "Y" ]]; then
        sudo "$(which shutdown)" -h now
    fi
}
function evaladvanced {
    local expression="$1"
    promptlog "$expression"
    eval "$expression"
}
function commandexists {
    local command="$1"
    if [ "$(command -v "$command")" ]; then
        echo true
    else
        echo false
    fi
}
function addalias {
    local name="$1"
    local command="$2"
    if [[ "--help" == "${name}" ]]||[[ "-h" == "${name}" ]]; then
        log "addalias NAME COMMAND"
        return
    fi
    delfilelines "$MY_ALIAS" "^alias $name="
    echo "alias $name=\"$command\"" | tee -a "${MY_ALIAS}" > /dev/null
}
alias editalias="nano '$MY_ALIAS'"
# TODO: Implement delalias
alias editprofile="nano '$MY_SHELL_PROFILE'"
alias editcustomprofile="nano '$MY_CUSTOM_SHELL_PROFILE'"
alias reloadprofile="source '$MY_SHELL_PROFILE'"
alias ver='lsb_release -a'
alias trash='mv --force -t ~/.local/share/Trash/files '
function createservice {
    local name="$1"
    local command="$2"
    local description="$3"
    local data="[Unit]
Description=$description

[Service]
ExecStart=$command

[Install]
WantedBy=multi-user.target
"
    echo "$data" | sudo tee "/etc/systemd/system/${name}.service" >/dev/null
    sudo systemctl start $name
    sudo systemctl enable $name
}
function deleteservice {
    local name="$1"
    if [ -f "$name" ]; then
        sudo systemctl stop $name
        sudo rm -rf "$name"
    fi
}
function restartexplorer {
    nautilus -q
}
# TODO: Implement setenv
# TODO: Implement deleteenv

# Features on linux only
alias update-menu-entries="sudo update-desktop-database"
alias cls='clear'
