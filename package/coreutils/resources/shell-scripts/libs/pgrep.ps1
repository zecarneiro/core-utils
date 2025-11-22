param(
    [Parameter(Mandatory=$true)]
    [string]$Pattern
)
$process = (Get-Process | Where-Object { $_.Name -match $Pattern })
$process | ForEach-Object {
    $id = $_.Id
    $name = $_.Name
    Write-Host "$id $name"
}