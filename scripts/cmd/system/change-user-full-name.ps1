# Author: José M. C. Noronha

$username = $env:username
Get-WmiObject Win32_UserAccount | ForEach-Object {
    $name = $_.Name
    $fullname = $_.FullName
    if ("$name" -eq "$username") {
        title-log "Change User Full name(Display name on start menu, etc)"
        Write-Host "Current user full name: $fullname"
        $newFullname = (Read-Host "Insert the new full name for the username '$username' (PRESS ENTER TO KEEP)")
        if (!([string]::IsNullOrEmpty($newFullname)) -and "$newFullname" -ne "$fullname") {
            sudopwsh Set-LocalUser -Name "$username" -FullName "'$newFullname'"
            ok-log "Change User Full name will be done when you logout or restart PC."
        }
    }
}