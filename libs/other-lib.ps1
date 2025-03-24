function install-visual-c-runtimes {
    $batScript = "$BIN_DIR\Visual-C-Runtimes-All-in-One\install_all.bat"
    if ((__show_install_message_question "Visual-C-Runtimes") -eq "y") {
        infolog "Install Visual-C-Runtimes"
        evaladvanced "sudopwsh $batScript"
    }
}

function install-powershellget {
    evaladvanced "sudopwsh Install-Module -Name PowerShellGet -Force"
}

function install-modules {
    evaladvanced "Install-Module -AllowClobber -Name scoop-completion -Scope CurrentUser -Force" # Project URL - https://github.com/Moeologist/scoop-completion"
    evaladvanced "Install-Module PSReadLine -Repository PSGallery -Scope CurrentUser -Force" # https://github.com/PowerShell/PSReadLine
}

function enable-sudo {
    $res = Read-Host "For Windows 11 only. Do you want to enable sudo? [y/N]"
    if ("$res" -eq "y" -or "$res" -eq "Y") {
        powershell -Command "Start-Process -Wait PowerShell -Verb RunAs -ArgumentList 'sudo.exe config --enable enable'"
    }
}

