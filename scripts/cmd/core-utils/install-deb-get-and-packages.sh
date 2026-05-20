#!/usr/bin/env bash

OPERATION_ARG="$1"
OS_VERSION_ARG="$2"

if [[ "${OS_VERSION_ARG}" == "26.04" ]]; then
    warn-log "deb-get can not be install on Ubuntu $OS_VERSION_ARG"
    exit 0
fi

function install_deb_get() {
    evalc "curl -sL https://raw.githubusercontent.com/wimpysworld/deb-get/main/deb-get | sudo -E bash -s install deb-get"
}

case "${OPERATION_ARG}" in
    manager) install_deb_get ;;
    *) warn-log "No Deb-Get operation to process." ;;
esac
