# Author: José M. C. Noronha and ChatGPT
# If User CANCEL the exit status code will be 130
param(
    [int] $timeToWait = 5
)

for($i = $timeToWait; $i -gt 0; $i--) {
    Write-Host "`r(Press ENTER to continue or others to CANCEL) Continue in ${i}... " -NoNewline
    $keyAvailable = $false
    try {
        $keyAvailable = [Console]::KeyAvailable
    } catch {
        # Host sem console interativo: segue contagem sem ler tecla
    }
    if ($keyAvailable) {
        if ([Console]::KeyAvailable) {
            $key = [Console]::ReadKey($true)
            if ($key.Key -ne [ConsoleKey]::Enter) {
                Write-Host ""
                exit 130 # SIGINT exit status
            }
            break
        }
    }
    Start-Sleep -Seconds 1
}
Write-Host ""
