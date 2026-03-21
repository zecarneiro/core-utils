# Author: José M. C. Noronha
[cmdletbinding()]
param (
    [string] $appId,
    [string]
    [parameter(ValueFromPipeline)]
    $message,
    [string] $icon,
    [Alias("h")]
    [switch] $help
)

if ($help) {
    Write-Host "$PSCommandPath -appId ID -title TITLE -message MESSAGE -icon ICON"
    exit 0
}
notify -appId "$appId" -title "Warning" -message "$message" -icon "$icon"