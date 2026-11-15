param(
    [string]$OPERATION_ARG = ""
)

switch ($OPERATION_ARG) {
    manager {
        evalc.exe "wsl.exe --update"
        sudo evalc.exe ". `"$PSScriptRoot\enable-system-features.ps1`""
    }
    Default { warn-log.exe "No WSL operation to process." }
}
