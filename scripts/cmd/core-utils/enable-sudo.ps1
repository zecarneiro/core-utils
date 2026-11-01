Write-Host "[INFO] Enabling sudo..."
powershell.exe -Command "Start-Process -Wait PowerShell -Verb RunAs -ArgumentList 'sudo.exe config --enable enable'"
