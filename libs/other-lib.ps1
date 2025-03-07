function install-visual-c-runtimes {
    $batScript = "$BIN_DIR\Visual-C-Runtimes-All-in-One\install_all.bat"
    if ((__show_install_message_question "Visual-C-Runtimes") -eq "y") {
        infolog "Install Visual-C-Runtimes"
        evaladvanced "sudopwsh $batScript"
    }
}