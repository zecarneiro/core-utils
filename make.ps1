$start = $false
$onlyProfile = $false

$SCRIPT_UTILS_DIR = ($PSScriptRoot)
$VERSION = (Get-Content "$SCRIPT_UTILS_DIR\version")
$SHELL_SCRIPT_DIR = "${SCRIPT_UTILS_DIR}\scripts"
$LIBS_DIR = "${SCRIPT_UTILS_DIR}\libs"
$BIN_DIR = "${SCRIPT_UTILS_DIR}\bin"
$IMAGES_DIR = "${SCRIPT_UTILS_DIR}\images"
$USER_OPTION_INSERTED_FILE = "$SCRIPT_UTILS_DIR\user-option-inserted"

if ($args[0] -eq "-s" -or $args[0] -eq "--start") {
    $start = $true
}
if ($args[1] -eq "--only-profile-shell") {
    $onlyProfile = $true
}

# ---------------------------------------------------------------------------- #
#                                    IMPORTS                                   #
# ---------------------------------------------------------------------------- #
Get-ChildItem ("${SHELL_SCRIPT_DIR}\*.ps1") | ForEach-Object {
    $fullname = $_.FullName
    . "$fullname"
}
Get-ChildItem ("${LIBS_DIR}\*.ps1") | ForEach-Object {
    $fullname = $_.FullName
    . "$fullname"
}

function manageUserOption {
    param([int] $option, [bool] $isGet)
    $doneRes = ""
    if ($isGet) {
        if ((filecontain "$USER_OPTION_INSERTED_FILE" "$option")) {
            $doneRes = " - Done"
        }
    } else {
        writefile -file "$USER_OPTION_INSERTED_FILE" -content "$option" -append
    }
    return "$doneRes"
}

function usage {
    Write-Host "Usage: make.ps1 [OPTIONS]... [STEP-VALUE]"
    Write-Host "OPTIONS:
     -s|--start`tProcess install and config by user
    "
}

function printMenu {
    echo "AAA"
    $(manageUserOption 3 $true)
    echo "BBB"
    Write-Host "1. Will$(manageUserOption 1 $true)
    - Enable Sudo
    - Set user bin dir
    - Install Scoop and Winget
    - Install PowershellGet Module
2. Will$(manageUserOption 2 $true)
    - Install Scoop and Winget packages
    - Install Powershell Modules
    - Install Visual-C-Runtimes
    - Install Features for WSL
3. Will$(manageUserOption 3 $true)
    - Create user powershell profile file
    - Install scripts profile
4. Will$(manageUserOption 4 $true)
    - Install Development packages. User will decide wich to install
    - Start all necessary configurations
5. Will(Optional)$(manageUserOption 5 $true)
    - Define/Change default system dirs. Like Documents, Images, etc
    - Change the user full display name
---
6. Exit"
}

function initProcess {
    $message = "Please, restart your terminal."
    __create_dirs
    if (!(__is_valid_home_dir)) {
        __show_rules_username
        __exit_script
    }
    for (;;) {
        $option = -1
        if (!($onlyProfile)) {
            printMenu
            $option = Read-Host "Insert an option"
        } else {
            $option = 3
        }
        switch ($option) {
            1 {
                __enable_sudo
                __set_user_bin_dir
                __install_scoop
                __install_winget
                __install_powershellget
                warnlog "$message"
                manageUserOption 1 $false
                __exit_script
            }
            2 {
                __install_scoop_packages
                __install_winget_packages
                __install_modules
                __install_visual_c_runtimes
                __install_features_for_wsl
                warnlog "After reboot, continue with option 3"
                manageUserOption 2 $false
                reboot
            }
            3 {
                __create_profile_file_powershell
                __install_profile_scripts
                if (!$onlyProfile) {
                    Set-Location "$SCRIPT_UTILS_DIR"
                    bash -c "./make.sh --start --only-profile-shell"
                    warnlog "$message"
                }
                manageUserOption 3 $false
                __exit_script
            }
            4 {
                __install_development_package
                __config_all
                manageUserOption 4 $false
                reboot
            }
            5 {
                __define_default_system_dir
                __change_user_full_name
                manageUserOption 5 $false
                reboot
            }
            6 { __exit_script }
            Default { Write-Host "WARN: Please, insert a valid option!" }
        }
    }
    
}

function main {
    if ($start) {
        initProcess
    } else {
        usage
    }
}
main
