param(
    [string]$Pattern,
    [Alias("h")]
    [switch]$help
)
function show_help() {
    Write-Host "pkill PATTERN"
}
if ($help) {
    show_help
    exit
}
if ((count-args $args) -le 0 -or (text-is-empty "$Pattern") -eq "true") {
    error-log "Invalid given Pattern: $Pattern"
    show_help
    exit
}
Get-Process $Pattern -ErrorAction SilentlyContinue | Stop-Process