Write-Host "[IMPORTANT] Disable Installer Dectection(Only continue if is the first time)"
$choice = Read-Host "Continue? [y/N]"
if ($choice -eq "y" -or $choice -eq "Y") {
    Write-Host [INFO] Enabling sudo...
    powershell.exe -Command "Start-Process -Wait PowerShell -Verb RunAs -ArgumentList 'sudo.exe config --enable enable'; pause"
    $cmdToRun = "Set-ItemProperty -Path `"HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`" -Name `"EnableInstallerDetection`" -Value 0"
    Write-Host ">>> $cmdToRun"
    sudo powershell.exe -Command "$cmdToRun"
    sudo powershell.exe -Command "Restart-Computer -Force"
}
