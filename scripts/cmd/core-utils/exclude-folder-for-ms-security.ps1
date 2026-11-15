$home_dir = "$env:USERPROFILE"

function Add-Dir() {
    param([string] $dir)
    Write-Host ">>> Add-MpPreference -ExclusionPath `"$dir`""
    Add-MpPreference -ExclusionPath "$dir"
}
Add-Dir "$home_dir\.local\bin"
Add-Dir "$home_dir\.local\opt"
Add-Dir "$home_dir\.local\share\coreutils"
