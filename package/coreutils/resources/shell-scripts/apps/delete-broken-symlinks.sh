#!/usr/bin/env bash
declare ARG="$*"

function show_help() {
  echo "delete-broken-symlinks DIR"
}

if [[ "${ARG}" == "-h" ]]||[[ "${ARG}" == "--help" ]]; then
  show_help
  exit 0
fi
if [[ -z "${ARG}" ]]; then
    ARG="."
fi
find "$ARG" -xtype l -exec rm {} \;