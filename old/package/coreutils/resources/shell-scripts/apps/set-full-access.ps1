param(
    [string]$file,
    [Alias("h")]
    [switch]$help
)
function show_help() {
    Write-Host "set-full-access FILE"
}
if ($help) {
    show_help
    exit
}
if ((count-args $args) -le 0 -or (text-is-empty "$file") -eq "true") {
    error-log "Invalid given Path"
    show_help
    exit
}
$acl = "$file" | Get-Acl -ErrorAction Stop
$rule = [security.accesscontrol.filesystemaccessrule]::new($Env:username, 'FullControl', 'Allow')
$acl.AddAccessRule($rule)
$acl | Set-Acl
