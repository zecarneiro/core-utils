$__PACKAGE_NAME__="coreutils"
$__MAIN_SCRIPT_DIR__ = ($PSScriptRoot)
$__PACKAGE_DIR__=(Resolve-Path "${__MAIN_SCRIPT_DIR__}\package")
$__BIN_DIR__=(Resolve-Path "${__MAIN_SCRIPT_DIR__}\bin")
$__CONFIG_DIR__=(Resolve-Path "${home}\.config")
$__SHELL_FILE__="$PROFILE"

# IMPORT LIBS FILES
Get-ChildItem -Path "$(Resolve-Path "${__MAIN_SCRIPT_DIR__}\libs")" -Filter *.ps1 -File | ForEach-Object {
  $lib_to_import = $_.FullName
  . "$lib_to_import"
}

function __create_dirs {
    $dirs = @("$__CONFIG_DIR__")
    Foreach ($dir in $dirs) {
        if (!(Test-Path -Path "$dir")) {
            __eval "New-Item -ItemType Directory -Force -Path `"$dir`" | Out-Null"
        }
    }
}

function __install_dependencies() {
    __create_dirs
    # Ensure the profile file exists
    if (-not (Test-Path -Path $__SHELL_FILE__)) {
        __eval "New-Item -ItemType File -Path `"$__SHELL_FILE__`" -Force | Out-Null"
    }
    if (__is_windows_so) {
        __set_user_bin_dir
        __eval ". `"$__SHELL_FILE__`"" $true
        . "$__SHELL_FILE__"
        __install_winget
        __install_scoop_and_packages
        __eval ". `"$__SHELL_FILE__`"" $true
        . "$__SHELL_FILE__"
        __install_visual_c_runtimes
        __install_features_for_wsl
    }
}

function __clean() {
    param([bool] $is_main_step = $true)
    $msg = "Cleanning"
    if ($is_main_step) {
        __print_title "$msg"
    } else {
        __print "$msg"
    }
    $directories = @("${__PACKAGE_DIR__}/build", "${__PACKAGE_DIR__}/coreutils.egg-info")
    foreach ($directory in $directories) {
        if (Test-Path -Path "$directory" -PathType Container) {
            __eval "Remove-Item -Path `"$directory`" -Recurse"
        }
    }
    __print "$msg. Done."
}

function __install() {
    param([bool] $is_main_step = $true)
    $msg = "Install ${__PACKAGE_NAME__}"
    if ($is_main_step) {
        __print_title "$msg"
    } else {
        __print "$msg"
    }
    if (__is_windows_so) {
        __uninstall $false
        Set-Location "$__PACKAGE_DIR__"
        __eval "pip-custom install ."
        Set-Location "$__MAIN_SCRIPT_DIR__"
    }
    __clean $false
    __eval "${__PACKAGE_NAME__}-postinstall"
    __print "$msg. Done."
}

function __uninstall() {
    param([bool] $is_main_step = $true)
    $msg = "Uninstall ${__PACKAGE_NAME__}"
    if ($is_main_step) {
        __print_title "$msg"
    } else {
        __print "$msg"
    } 
    __eval "${__PACKAGE_NAME__}-preuninstall"
    if (__is_windows_so) {
        __eval "pip-custom uninstall ${__PACKAGE_NAME__} --yes"
    }
    __print "$msg. Done."
}

function __show_help {
    Write-Host "Usage: Main.ps1 [OPTION]"
    Write-Host
    Write-Host "Options:"
    Write-Host "  -d, --install-dependencies   Install Python and required dependencies."
    Write-Host "                               After installation, restart your terminal."
    Write-Host
    Write-Host "  -i, --install                Uninstall any previous installation and reinstall the project."
    Write-Host
    Write-Host "  -u, --uninstall              Uninstall the project and remove related files."
    Write-Host
    Write-Host "  -c, --clean                  Clean temporary files and caches."
    Write-Host
    Write-Host "  -h, --help                   Show this help message and exit."
    Write-Host
}

function main {
    param([string]$Option)

    switch ($Option) {
        "-d" { Main "--install-dependencies" }
        "--install-dependencies" {
            __install_dependencies
            Write-Host "[INFO] Please, restart your terminal!"
        }
        "-i" { Main "--install" }
        "--install" {
            __install
        }

        "-u" { Main "--uninstall" }
        "--uninstall" {
            __uninstall
        }

        "-c" { Main "--clean" }
        "--clean" {
            __clean
        }

        "-h" { Main "--help" }
        "--help" {
            __show_help
        }

        Default {
            Write-Host "[ERROR] Invalid option: $Option"
            Write-Host "Use 'Main.ps1 --help' to see available options."
        }
    }
}

# Call the main function with the first argument
main $args[0]

