param(
    [string]$OPERATION_ARG = ""
)

$feature = Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All
if ($feature.State -ne "Enabled") {
    Write-Host "=== Hyper-V is not enabled. Enabling ==="
    Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V-All -All -NoRestart
    Write-Host "[OK] Hyper-V enabled. A restart is required to take efect."
} else {
    Write-Host "[OK] Hyper-V is already enabled."
}
