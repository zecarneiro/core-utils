function install-scoop {
    if (!(commandexists scoop)) {
        Write-Host "INFO: Install Scoop ..."
        Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
    }
}

function install-scoop-packages {
    evaladvanced "scoop install main/coreutils"
    evaladvanced "scoop install main/git"
    evaladvanced "scoop install main/vim"
    evaladvanced "scoop install main/nano"
    evaladvanced "scoop install main/curl"
    evaladvanced "scoop install main/grep"
    evaladvanced "scoop install main/sed"
    evaladvanced "scoop install main/which"
    evaladvanced "scoop install main/dos2unix"
    evaladvanced "scoop install main/7zip"
    evaladvanced "scoop install main/gsudo"

    evaladvanced "scoop bucket add extras"
    evaladvanced "scoop install extras/okular"
    evaladvanced "scoop install extras/git-credential-manager"
    # evaladvanced "scoop install extras/scoop-completion" # Project URL - https://github.com/Moeologist/scoop-completion
    evaladvanced "scoop bucket add alkuzad_unxutils-separated https://github.com/alkuzad/unxutils-separated"
    evaladvanced "scoop install alkuzad_unxutils-separated/unxutils-xargs"
}

function config-scoop-all {
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
    log "Step 1: Open 7-Zip File Manager(As Admin) by typing 7-Zip in Start menu or Start screen and then pressing Enter key."
    log "Step 2: Next, navigate to Tools menu and then click Options to open Options dialog."
    log "Step 3: Here, under 7-Zip tab, make sure that Integrate 7-Zip to shell context menu option is selected. If not, please select the option and then click Apply button. You might need to reboot your PC or restart Windows Explorer to see 7-Zip in the context menu."
    pause
}