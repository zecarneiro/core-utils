if ([string]::IsNullOrWhiteSpace($args)) {
    error-log "No args passed"
    exit
}
Start-Process powershell.exe -verb runAs -Args "$args; pause"