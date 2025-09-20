#!/usr/bin/env bash
# Author: José M. C. Noronha

alias zenity="zenity 2>/dev/null"

function notify {
    local appId="$1"
    local title="$2"
    local message="$3"
    local icon="$4"
    if [[ -n "${icon}" ]]; then
        notify-send -a "$appId" -i "$icon" -t 60000 "$title" "$message"
    else
        notify-send -a "$appId" -t 60000 "$title" "$message"
    fi
}
function oknotify {
    local appId="$1"
    local message="$2"
    local icon="$3"
    notify "$appId" "Success" "$message" "$icon"
}
function infonotify {
    local appId="$1"
    local message="$2"
    local icon="$3"
    notify "$appId" "Information" "$message" "$icon"
}
function warnnotify {
    local appId="$1"
    local message="$2"
    local icon="$3"
    notify "$appId" "Warning" "$message" "$icon"
}
function errornotify {
    local appId="$1"
    local message="$2"
    local icon="$3"
    notify "$appId" "Error" "$message" "$icon"
}
function selectfiledialog {
    local selectedFile=$(zenity --file-selection --modal)
    echo "{ \"selected\": \"$selectedFile\" }"
}
