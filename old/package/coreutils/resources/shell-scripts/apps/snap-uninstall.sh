#!/usr/bin/env bash
declare ARG="$*"

function show_help() {
  echo "snap-uninstall APP"
}

if [[ "${ARG}" == "-h" ]]||[[ "${ARG}" == "--help" ]]; then
  show_help
  exit 0
fi
if [[ $(count-args $ARG) -le 0 ]]||[[ $(text-is-empty "$ARG") == "true" ]]; then
  error-log "Invalid given APP"
  show_help
  exit 0
fi

declare config_dir="$HOME/snap/$ARG"
declare config_system_dir="/snap/$ARG"
evalc "sudo snap remove --purge '$ARG'"
evalc "snap saved"
read -p "Insert the number on the line of App(ENTER TO SKIP): " userInput
if [[ -n "${userInput}" ]]; then
    evalc "sudo snap forget ${userInput}"
fi
if [ -d "$config_dir" ]; then
    evalc "rm -rf \"$config_dir\""
fi
if [ -d "$config_system_dir" ]; then
    evalc "sudo rm -rf \"$config_system_dir\""
fi
