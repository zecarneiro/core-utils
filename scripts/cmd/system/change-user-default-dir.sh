#!/usr/bin/env bash
# shellcheck disable=SC2162

declare -A user_dirs
read -e -p "Insert DOWNLOAD (Or ENTER to cancel) " result
if [ -d "$result" ]; then
    user_dirs[DOWNLOAD]="$result"
fi
read -e -p "Insert TEMPLATES (Or ENTER to cancel) " result
if [ -d "$result" ]; then
    user_dirs[TEMPLATES]="$result"
fi
read -e -p "Insert DOCUMENTS (Or ENTER to cancel) " result
if [ -d "$result" ]; then
    user_dirs[DOCUMENTS]="$result"
fi
read -e -p "Insert MUSIC (Or ENTER to cancel) " result
if [ -d "$result" ]; then
    user_dirs[MUSIC]="$result"
fi
read -e -p "Insert PICTURES (Or ENTER to cancel) " result
if [ -d "$result" ]; then
    user_dirs[PICTURES]="$result"
fi
read -e -p "Insert VIDEOS (Or ENTER to cancel) " result
if [ -d "$result" ]; then
    user_dirs[VIDEOS]="$result"
fi
for key in "${!user_dirs[@]}"; do
    info-log "Change $key"
    xdg-user-dirs-update --set $key "${user_dirs[$key]}"
done