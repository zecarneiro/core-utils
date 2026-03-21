# Author: José M. C. Noronha
[cmdletbinding()]
param (
    [string] $message,
    [Alias("h")]
    [switch] $help
)

if ($help) {
    Write-Host "$PSCommandPath -appId ID -title TITLE -message MESSAGE -icon ICON"
    exit 0
}
Write-Host "$message"
Add-Type -AssemblyName System.Windows.Forms
$browser = New-Object System.Windows.Forms.FolderBrowserDialog
$null = $browser.ShowDialog()
$selectedDir = $browser.SelectedPath
return @{ selected="$selectedDir"; }