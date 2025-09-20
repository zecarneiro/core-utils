#!/usr/bin/env bash
declare ARG="$*"

function show_help() {
  echo "pgrep PATTERN"
}

if [[ "${ARG}" == "-h" ]]||[[ "${ARG}" == "--help" ]]; then
  show_help
  exit 0
fi
if [[ $(count-args $ARG) -le 0 ]]||[[ $(text-is-empty "$ARG") == "true" ]]; then
  error-log "Invalid given Pattern"
  show_help
  exit 0
fi
$(which "pgrep") -l "$ARG"
