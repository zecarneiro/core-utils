#!/usr/bin/env bash
# Author: José M. C. Noronha

function fileexists {
    local file="$1"
    if [ -f "${file}" ]; then
        echo true
    else
        echo false
    fi
}
function fileextension {
    local file="$1"
    if [[ $(echo "$file" | grep -o "\." | wc -l) -gt 0 ]]; then
        echo "${file##*.}"
    fi
}
function filename {
    local file="$1"
    if [[ -n "${file}" ]]; then
        echo "${file%.*}"
    fi
}
function writefile {
    local file="$1"
    local content="$2"
    local append="$3"
    if [[ "--help" == "${file}" ]]||[[ "-h" == "${file}" ]]; then
        log "writefile FILE CONTENT [ |APPEND]"
    elif
        if [[ ! -f "$file" ]]||[[ $append == false ]]; then
            echo -e "$content" | tee "$file" >/dev/null
        else
            echo -e "$content" | tee -a "$file" >/dev/null
        fi
    fi
}
function delfilelines {
    local file="$1"
    local match="$2"
    if [[ "--help" == "${file}" ]]||[[ "-h" == "${file}" ]]; then
        log "delfilelines FILE MATCH"
        return
    fi
    if [ "$(filecontain "$file" "$match")" == true ]; then
        local tempfile="${file}.tmp"
        cat "$file" | grep -v "$match" > "${tempfile}"
        mv "$tempfile" "$file"
    fi
}
function deletefile {
    local file="$1"
    if [[ $(fileexists "$file") == true ]]; then
        rm -rf "$file"
        infolog "Deleted file: $file"
    fi
}
alias countfiles="find . -type f | wc -l"
alias findfile="find . | grep "
alias movefilestoparent="find . -mindepth 2 -type f -print -exec mv {} . \;"
alias lf="find \".\" -maxdepth 1 -type f"
function filecontain {
    local file="$1"
    local match="$2"
    if [[ "--help" == "${file}" ]]||[[ "-h" == "${file}" ]]; then
        log "filecontain FILE MATCH"
        return
    fi
    if [ "$(cat "$file" | grep -c "$match")" -ne 0 ]; then
        echo true
    else
        echo false
    fi
}
function openmarkdown {
    local file="${1}"
    if [[ $(fileexists "$file") == true ]]; then
        if [ "$(commandexists "ghostwriter")" == true ]; then
            ghostwriter "$file" &
        else
            cat "$file"
        fi
    fi
}
function openimage {
    local file="${1}"
    if [[ $(fileexists "$file") == true ]]; then
        if [ "$(commandexists "eog")" == true ]; then
            eog "$file" &
        else
            errorlog "Can not open given image. Command not found: eog"
        fi
    fi
}
