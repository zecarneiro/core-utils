param(
    [string]$Pattern
)
Get-Process $Pattern -ErrorAction SilentlyContinue | Stop-Process