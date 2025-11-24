#!/usr/bin/env bash
declare __PACKAGE_NAME__="coreutils"
declare __MAIN_SCRIPT_DIR__="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
declare __PACKAGE_DIR__="${__MAIN_SCRIPT_DIR__}/package"
declare __CONFIG_DIR__="${HOME}/.config"
declare __SCRIPTS_DIR__="${__MAIN_SCRIPT_DIR__}/scripts"
declare __SHELL_FILE__="${HOME}/.bashrc"

# IMPORT LIBS FILES
for lib_to_import in "$__MAIN_SCRIPT_DIR__"/libs/*.sh; do
  # shellcheck disable=SC1090
  . "$lib_to_import"
done

function __create_dirs {
    local dirs=("$__CONFIG_DIR__")
    for dir in "${dirs[@]}"; do
        if [[ ! -d "${dir}" ]]; then
            __eval "mkdir '$dir'"
        fi
    done
}

function __install_dependencies() {
    __create_dirs
    if [[ $(__is_linux_so) == true ]]; then
        __install_scripts
        __set_user_bin_dir

        __eval ". '$__SHELL_FILE__'" true
        . "$__SHELL_FILE__"
        
        __install_apt_and_packages
        __install_flatpak_and_packages
        __install_snap_and_packages
        __install_appimage_and_packages
        __install_pacstall_and_packages
        __install_deb_get_and_packages

        __eval ". '$__SHELL_FILE__'" true
        . "$__SHELL_FILE__"

        __install_rust
    fi
}

function __clean() {
    local is_main_step="${1}"
    local msg="Cleanning"
    if [[ "${is_main_step}" == true ]]; then
        __print_title "$msg"
    else
        __print "$msg"
    fi
    local -a directories=("${__PACKAGE_DIR__}/build" "${__PACKAGE_DIR__}/coreutils.egg-info")
    for directory in "${directories[@]}"; do
        if [ -d "${directory}" ]; then
            __eval "rm -r '${directory}'"
        fi
    done
    __print "$msg. Done."  
}

function __build() {
    pushd .
    cd "$__PACKAGE_DIR__" || exit 1
    __eval "ruff check ." || exit 1
    __eval "pyright" || exit 1
    popd
}

function __install() {
    local is_main_step="${1}"
    __build
    local msg="Install ${__PACKAGE_NAME__}"
    if [[ "${is_main_step}" == true ]]; then
        __print_title "$msg"
    else
        __print "$msg"
    fi
    if [[ $(__is_linux_so) == true ]]; then
        __uninstall false
        (
            cd "$__PACKAGE_DIR__" || exit
            __eval "pipc install ."
        )
    fi
    __clean false
    __eval "${__PACKAGE_NAME__}-postinstall"
    __eval "pwsh -Command ${__PACKAGE_NAME__}-postinstall"
    __print "$msg. Done."
}

function __uninstall() {
    local is_main_step="${1}"
    local msg="Uninstall ${__PACKAGE_NAME__}"
    if [[ "${is_main_step}" == true ]]; then
        __print_title "$msg"
    else
        __print "$msg"
    fi
    __eval "${__PACKAGE_NAME__}-preuninstall"
    __eval "pwsh -Command ${__PACKAGE_NAME__}-preuninstall"
    if [[ $(__is_linux_so) == true ]]; then
        __eval "pipc uninstall ${__PACKAGE_NAME__} --yes"
    fi
    __print "$msg. Done."
}

function main {
    case "${1}" in
        -d|--install-dependencies)
            __install_dependencies
            echo "[INFO] Please, restart your terminal!"
        ;;
        -i|--install)
            __install true
        ;;
        -u|--uninstall)
            __uninstall true
        ;;
        -c|--clean)
            __clean true
        ;;
        -h|--help)
            echo "Usage: $0 [OPTION]"
            echo
            echo "Options:"
            echo "  -d, --install-dependencies   Install Python and required dependencies."
            echo "                               After installation, restart your terminal."
            echo
            echo "  -i, --install                Uninstall any previous installation and reinstall the project."
            echo
            echo "  -u, --uninstall              Uninstall the project and remove related files."
            echo
            echo "  -c, --clean                  Clean temporary files and caches."
            echo
            echo "  -h, --help                   Show this help message and exit."
            echo
        ;;
        *)
            echo "[ERROR] Invalid option: ${1:-none}"
            echo "Use '$0 --help' to see available options."
        ;;
    esac
}
main "$@"
