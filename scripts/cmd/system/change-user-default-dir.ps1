param(
    [string] $shellProfileDir
)

if ((is-admin) -eq "false") {
    sudopwsh "$PSScriptRoot\change-user-default-dir.ps1 `"$shellProfileDir`""
    exit 0
}

$isDocumentsChange = $false
$newShellProfileDir = ""
$userDirs = @{}
$isSetDirs = $false
$result=$(select-folder-dialog "Insert DOWNLOAD (Or cancel)")
$result = $result.selected
if (! [string]::IsNullOrEmpty($result)) {
    $userDirs.Add("{374DE290-123F-4565-9164-39C4925E467B}", "$result")
}
$result=$(select-folder-dialog "Insert DOCUMENTS (Or cancel)")
$result = $result.selected
if (! [string]::IsNullOrEmpty($result)) {
    $userDirs.Add("Personal", "$result")
    $newShellProfileDir = "$result\$(basename "$shellProfileDir")"
    $isDocumentsChange = $true
}
$result=$(select-folder-dialog "Insert MUSIC (Or cancel)")
$result = $result.selected
if (! [string]::IsNullOrEmpty($result)) {
    $userDirs.Add("My Music", "$result")
}
$result=$(select-folder-dialog "Insert PICTURES (Or cancel)")
$result = $result.selected
if (! [string]::IsNullOrEmpty($result)) {
    $userDirs.Add("My Pictures", "$result")
}
$result=$(select-folder-dialog "Insert VIDEOS (Or cancel)")
$result = $result.selected
if (! [string]::IsNullOrEmpty($result)) {
    $userDirs.Add("My Video", "$result")
}
foreach ($userDir in $userDirs.GetEnumerator()) {
    $isSetDirs=$true
    info-log "Change $($userDir.Name)"
    reg add "HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders" /f /v "$($userDir.Name)" /t REG_SZ /d "$($userDir.Value)"
}
if ($isDocumentsChange) {
    if ((direxists "$shellProfileDir")) {
        cpdir -s "$shellProfileDir" -d "$newShellProfileDir" -f
    }
}
if ($isSetDirs){
    restart-explorer
}
