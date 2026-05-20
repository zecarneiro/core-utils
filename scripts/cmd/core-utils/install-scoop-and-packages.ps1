param(
    [string]$OPERATION_ARG = "",
    [string]$SCOOP_DIR_ARG = ""
)

function Install-Scoop {
    evalc "Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression"
    evalc "Install-Module -AllowClobber -Name scoop-completion -Scope CurrentUser -SkipPublisherCheck"
}

function Install-Packages {
    $cmdList = @(
        "install main/7zip",
        "install main/git",
        "install main/openssl",
        "install main/vim",
        "install main/nano",
        "install main/curl",
        "install main/sed",
        "install main/clink",
        "install main/dos2unix",
        "install main/fzf",
        "bucket add extras",
        "install extras/psfzf",
        "install extras/psreadline", # https://github.com/PowerShell/PSReadLine
        "install extras/git-credential-manager",
        "install https://github.com/c3er/mdview/releases/latest/download/mdview.json"
    )
    foreach ($cmd in $cmdList) {
        evalc "scoop $cmd"
    }
}

function ConfigPackages {
    evalc "sudo cmd.exe /C ${SCOOP_DIR_ARG}\apps\7zip\current\install-context.reg"
    evalc "sudo cmd.exe /C ${SCOOP_DIR_ARG}\apps\git\current\install-associations.reg"
    evalc "clink set clink.logo none"
}

switch ($OPERATION_ARG) {
    manager { Install-Scoop }
    packages { Install-Packages }
    config { ConfigPackages }
    Default { warn-log "No Scoop operation to process." }
}
