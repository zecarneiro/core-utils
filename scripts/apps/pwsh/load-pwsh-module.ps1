param(
    [string]$name,
    [Alias("h")]
    [switch]$help
)
function show_help() {
    Write-Host "load-pwsh-module NAME"
}
if ($help) {
    show_help
    exit 0
}
if ([string]::IsNullOrWhiteSpace("$name")) {
    error-log "Invalid given module Name: $name"
    exit 1
}
Import-Module $name -ErrorAction SilentlyContinue
if (-not $?) {
    error-log "Import Module: $name"
}