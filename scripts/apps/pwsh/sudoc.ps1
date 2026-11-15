if ([string]::IsNullOrWhiteSpace($args))
{
    error-log "No args passed"
    exit 1
}
Write-Host ">>> " -ForegroundColor DarkGray -NoNewline; Write-Host "$args"
sudopwsh "$args"
