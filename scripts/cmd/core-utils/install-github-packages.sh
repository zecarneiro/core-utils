#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_packages() {
    info-log "Nothing to process..."
}

case "${OPERATION_ARG}" in
    packages) install_packages ;;
    *) warn-log "No GitHub Packages operation to process." ;;
esac
