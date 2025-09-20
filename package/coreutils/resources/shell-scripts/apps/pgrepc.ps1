param(
    [string]$Pattern,
    [Alias("h")]
    [switch]$help
)
function show_help() {
    Write-Host "pgrep PATTERN"
}
if ($help) {
    show_help
    exit
}
if ((count-args $args) -le 0 -or (text-is-empty "$Pattern") -eq "true") {
    error-log "Invalid given Pattern"
    show_help
    exit
}
$process = (Get-Process | Where-Object { $_.Name -match $Pattern })
$process | ForEach-Object {
    $id = $_.Id
    $name = $_.Name
    Write-Host "$id $name"
}