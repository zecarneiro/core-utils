#!/usr/bin/env bash
# Author: José M. C. Noronha

declare NO_COLOR='\033[0m'
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
    local RED_COLOR='\033[0;31m'
    log "[${RED_COLOR}ERROR${NO_COLOR}] ${message}" "$keep_line"
}
function infolog {
    local message="$1"
    local keep_line=$2
    local BLUE_COLOR='\033[0;34m'
    log "[${BLUE_COLOR}INFO${NO_COLOR}] ${message}" "$keep_line"
}
function debuglog {
    local message="$1"
    local keep_line=$2
    local YELLOW_COLOR='\033[0;33m'
    log "[DEBUG] ${message}" "$keep_line"
}
function warnlog {
    local message="$1"
    local keep_line=$2
    local YELLOW_COLOR='\033[0;33m'
    log "[${YELLOW_COLOR}WARN${NO_COLOR}] ${message}" "$keep_line"
}
function oklog {
    local message="$1"
    local keep_line=$2
    local GREEN_COLOR='\033[0;32m'
    log "[${GREEN_COLOR}OK${NO_COLOR}] ${message}" "$keep_line"
}
function promptlog {
    local message="$1"
    local DARKGRAY_COLOR='\033[1;30m'
    declare BOLD_FOR_COLOR='\e[0m'
    log "${DARKGRAY_COLOR}>>>${BOLD_FOR_COLOR}${NO_COLOR} ${BOLD_FOR_COLOR}${message}${NO_COLOR}"
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
