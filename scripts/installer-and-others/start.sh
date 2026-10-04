#!/bin/bash
operation="$1"
coreutil_rootdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$coreutil_rootdir" || exit 1

case "$operation" in
    "install")
        chmod +x "${coreutil_rootdir}/installers/install"
        eval "\"${coreutil_rootdir}/installers/install\" \"${coreutil_rootdir}\""
        ;;
    "uninstall")
        chmod +x "${coreutil_rootdir}/installers/uninstall"
        eval "\"${coreutil_rootdir}/installers/uninstall\""
        ;;
    *)
        echo "Uso: start.sh [install|uninstall]"
        exit 1
        ;;
esac