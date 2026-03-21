package main

var linuxTermplate = `[Desktop Entry]
Version=1.0
Type=Application
Name=%s
Exec=%s %s
Icon=%s
Categories=%s
Comment=%s
Terminal=%s
X-AppStream-Ignore=true
`

var windowsTemplate = `
$WshShell = New-Object -ComObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("%s")
$Shortcut.TargetPath = "%s"
$Shortcut.WindowStyle = %d
%s
$Shortcut.Save()
`

var windowsAdminTemplate = `
$path = "%s"
$bytes = [System.IO.File]::ReadAllBytes($path)
$bytes[0x15] = $bytes[0x15] -bor 0x20
[System.IO.File]::WriteAllBytes($path, $bytes)
`
