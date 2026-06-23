#!/usr/bin/env bash

OPERATION_ARG="$1"

function install_deb_get() {
    evalc "curl -sL https://raw.githubusercontent.com/wimpysworld/deb-get/main/deb-get | sudo -E bash -s install deb-get"
    source "$(SHELL_PROFILE_SCRIPT)"
    evalc "deb-get update"
}

case "${OPERATION_ARG}" in
    manager) install_deb_get ;;
    *) warn-log "No Deb-Get operation to process." ;;
esac
