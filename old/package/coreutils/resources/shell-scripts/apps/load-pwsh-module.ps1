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
    exit
}
if ((count-args $args) -le 0 -or (text-is-empty "$name") -eq "true") {
    error-log "Invalid given Name: $name"
    show_help
    exit
}
Import-Module $name -ErrorAction SilentlyContinue
if (-not $?) {
    error-log "Import Module: $name"
}