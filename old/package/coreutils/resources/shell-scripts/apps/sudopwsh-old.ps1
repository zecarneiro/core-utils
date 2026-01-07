if ((text-is-empty "$args") -eq "true") {
    error-log "No args passed"
    exit
}
Start-Process powershell.exe -verb runAs -Args "$args; pause"