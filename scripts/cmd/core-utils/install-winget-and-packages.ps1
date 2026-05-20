param(
    [string]$OPERATION_ARG = "",
    [string]$WingetInstallerScript = ""
)

function Install-Winget {
    sudo evalc ". `"$WingetInstallerScript`""
}

function Install-Packages {
    $idList = @(
        "topgrade-rs.topgrade"
        "Microsoft.Coreutils"
    )
    foreach ($id in $idList) {
        winget-install "$id"
    }
}

switch ($OPERATION_ARG) {
    manager { Install-Winget }
    packages { Install-Packages }
    Default { warn-log "No Winget operation to process." }
}
