# ---------------------------------------------------------------------------- #
#                                    WINGET                                    #
# ---------------------------------------------------------------------------- #
function __install_winget_packages {
    if ((__show_install_message_question "Do you to install Winget packages") -eq "y") {
        evaladvanced "winget install --id=chrisant996.Clink --accept-source-agreements --accept-package-agreements"
    }
}

# ---------------------------------------------------------------------------- #
#                                     SCOOP                                    #
# ---------------------------------------------------------------------------- #


function __install_scoop_packages {
    $scoopDir = "$home\scoop"
    if ((__show_install_message_question "Do you to install Scoop packages") -eq "y") {
        evaladvanced "scoop install main/uutils-coreutils"
        evaladvanced "scoop install main/git"
        evaladvanced "scoop install main/vim"
        evaladvanced "scoop install main/nano"
        evaladvanced "scoop install main/curl"
        evaladvanced "scoop install main/grep"
        evaladvanced "scoop install main/sed"
        evaladvanced "scoop install main/dos2unix"
        evaladvanced "scoop install main/7zip"

        evaladvanced "scoop bucket add extras"
        log "If you get error, please, open new terminal and run: scoop bucket add extras"
        pause
        
        evaladvanced "scoop install extras/git-credential-manager"
        evaladvanced "scoop bucket add alkuzad_unxutils-separated https://github.com/alkuzad/unxutils-separated"
        log "If you get error, please, open new terminal and run: scoop bucket add alkuzad_unxutils-separated https://github.com/alkuzad/unxutils-separated"
        pause
        evaladvanced "scoop install alkuzad_unxutils-separated/unxutils-xargs"

        # Markdown apps
        # evaladvanced "scoop install extras/marktext" # https://github.com/marktext/marktext
        evaladvanced "scoop install https://github.com/c3er/mdview/releases/latest/download/mdview.json" # https://github.com/c3er/mdview

        ## Config ##
        evaladvanced "gsudo config CacheMode auto"
    
        # Delete all system alias
        delalias "cp"
        delalias "cat"
        delalias "mkdir"
        delalias "ls"
        delalias "mv"
        delalias "ps"
        delalias "rm"
        delalias "rmdir"
        delalias "sleep"
        delalias "sort"
        delalias "tee"
        delalias "curl"
        delalias "grep"
        delalias "sed"

        # Docs
        titlelog "Integrate 7zip on context menu"
        evaladvanced "$scoopDir\apps\7zip\current\install-context.reg"
        pause
    }
}

# ---------------------------------------------------------------------------- #
#                               POWERSHELL MODULE                              #
# ---------------------------------------------------------------------------- #
function __install_powershell_modules {
    infolog "No powershell modules to process"
}

# ---------------------------------------------------------------------------- #
#                                    OTHERS                                    #
# ---------------------------------------------------------------------------- #
