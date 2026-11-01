# Certifique-se de que está a correr como Administrador
param(
    [switch] $help
)
if ($help) {
    Write-Host "Will be disable Device Guard and Smart App Control"
    Write-Host "Will be enable Installer Detection"
    exit 0
}
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    sudo powershell.exe -Command "$PSScriptRoot\windows-policies-processor.ps1"
    exit 0
}

Write-Host "[INFO] Init disable Device Guard / Smart App Control..."
reg add "HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\CI\Policy" /v VerifiedAndReputablePolicyState /t REG_DWORD /d 0 /f
Write-Host "[OK] Device Guard / Smart App Control was disabled."

Write-Host "[INFO] Init enable Installer Detection..."
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "EnableInstallerDetection" -Value 0
Write-Host "[OK] Installer Detection was enabled."

Write-Host "==============="
Write-Host "[WARN] Please, restart your PC to policies changes take effects."
