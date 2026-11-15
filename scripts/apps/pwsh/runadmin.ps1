if ([string]::IsNullOrWhiteSpace($args)) {
    error-log "No args passed"
    exit
}
Write-Host ">>> " -ForegroundColor DarkGray -NoNewline; Write-Host "$args"
Start-Process powershell.exe -ArgumentList "-NoProfile -Command `"$args; Read-Host 'Press ENTER to close...'`"" -Verb runAs -Wait
