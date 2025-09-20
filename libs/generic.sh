#!/usr/bin/env bash

function __eval() {
    local command="$1"
    local showOnly="${2:-false}"
    echo ">>> $command"
    if [[ "${showOnly}" == "false" ]]; then
        eval "$command"
    fi    
}

function __print() {
    echo "==== $1 ===="
}

function __print_title() {
    local msg="$*"
    local border
    border=$(printf '=%.0s' $(seq 1 $((${#msg} + 4))))
    echo "$border"
    echo "= $msg ="
    echo "$border"
}

function __is_linux_so() {
    if [ "$(uname)" = "Linux" ]; then
        echo true
    else
        echo false
    fi
}

function __file_exist_with_sudo() {
    sudo test -f "$1" && echo true|| echo false
}