function __eval() {
    param([string] $command, [bool] $showOnly = $false)
    Write-Host ">>> $command"
    if (!$showOnly) {
        Invoke-Expression "$command"
    }
}

function __print() {
    param([string] $data)
    Write-Host "==== $data ===="
}

function __print_title() {
    param([string] $data)
    $border = "=" * ($data.Length + 4)
    Write-Host $border
    Write-Host ("= " + $data + " =")
    Write-Host $border
}

function __is_windows_so() {
    if ([System.Environment]::OSVersion.Platform -eq "Win32NT") {
        return $true
    } else {
        return $false
    }
}