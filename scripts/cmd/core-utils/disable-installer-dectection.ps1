Write-Host "[IMPORTANT] Disable Installer Dectection(Only continue if is the first time)"
$choice = Read-Host "Continue? [y/N]"
if ($choice -eq "y" -or $choice -eq "Y") {
    $cmdToRun = "Set-ItemProperty -Path `"HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`" -Name `"EnableInstallerDetection`" -Value 0"
    Write-Host ">>> $cmdToRun"
    sudo powershell.exe -Command "$cmdToRun"
    sudo powershell.exe -Command "Restart-Computer -Force"
}
