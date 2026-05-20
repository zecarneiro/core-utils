param(
    [string]$OPERATION_ARG = ""
)

switch ($OPERATION_ARG) {
    manager { evalc "wsl.exe --install" }
    Default { warn-log "No WSL operation to process." }
}