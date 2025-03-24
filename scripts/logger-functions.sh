#!/usr/bin/env bash
# Author: José M. C. Noronha

function log {
    local message="$1"
    local keep_line=$2
    if [[ -n "${keep_line}" ]]&&[[ "${keep_line}" == "1" ]]; then
        echo -n "${message}" >&2
    else
        echo -e "${message}" >&2
    fi
}
function errorlog {
    local message="$1"
    local keep_line=$2
    log "[${RedColor}ERROR${ResetColor}] ${message}" "$keep_line"
}
function infolog {
    local message="$1"
    local keep_line=$2
    log "[${BlueColor}INFO${ResetColor}] ${message}" "$keep_line"
}
function debuglog {
    local message="$1"
    local keep_line=$2
    log "[DEBUG] ${message}" "$keep_line"
}
function warnlog {
    local message="$1"
    local keep_line=$2
    log "[${YellowColor}WARN${ResetColor}] ${message}" "$keep_line"
}
function oklog {
    local message="$1"
    local keep_line=$2
    log "[${GreenColor}OK${ResetColor}] ${message}" "$keep_line"
}
function promptlog {
    local message="$1"
    log "${DarkGrayColor}>>>${BoldColor}${ResetColor} ${BoldColor}${message}${ResetColor}"
}
function titlelog {
    local message="$1"
    local message_len=${#message}
    local separator=""
    for (( i=1; i<=message_len+8; i++ )); do
        separator="#${separator}"
    done
    log "$separator"
    log "##  $message  ##"
    log "$separator"
}
function headerlog {
    local message="$1"
    log "##  $message  ##"
}
function separatorlog {
    local length=$1
    local message="# "
    if [[ -z "${length}" ]]||[[ $length -lt 5 ]]; then
        length=6
    fi
    for (( i=1; i<=length-4; i++ )); do
        message="${message}-"
    done
    message="$message #"
    log "$message"
}
