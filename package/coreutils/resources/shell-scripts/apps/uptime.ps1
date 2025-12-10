if ($PSVersionTable.PSVersion.Major -eq 5 ) {
    Get-WmiObject win32_operatingsystem | Select-Object @{EXPRESSION = { $_.ConverttoDateTime($_.lastbootuptime) } } | Format-Table -HideTableHeaders
} else {
    net statistics workstation | Select-String "since" | foreach-object { $_.ToString().Replace('Statistics since ', '') }
}