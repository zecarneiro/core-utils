if ([string]::IsNullOrWhiteSpace($args)) {
    error-log "No args passed"
    exit
}
# OLD VERSION: Start-Process powershell.exe -verb runAs -Args "$args; pause"
sudo powershell.exe -Command "$args"
