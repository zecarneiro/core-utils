param(
    [Parameter(Mandatory=$true)]
    [string]$Pattern
)
Get-Process $Pattern -ErrorAction SilentlyContinue | Stop-Process