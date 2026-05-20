param(
    [string]$OPERATION_ARG = "",
    [bool]$FORCE_ARG = $false,
    [string]$SCRIPTS_APPS_DIR_ARG="",
    [string]$SCOOP_DIR_ARG=""
)
$SCRIPT_FORCE_ARG = ""
if ($FORCE_ARG) {
    $SCRIPT_FORCE_ARG = "-f"
}

function Install-ScriptsApps {
    Get-ChildItem -Path "$SCRIPTS_APPS_DIR_ARG" -File | ForEach-Object {
        $scriptApp = $_.FullName
        script-manager-cu install "$scriptApp" $SCRIPT_FORCE_ARG
    }
}

function Install-Alias {
    $aliasMap = @{
        "pwsh" = "powershell.exe"
        "now" = "date"
        "bash" = "& '$SCOOP_DIR_ARG\apps\git\current\bin\bash.exe'" # I put & because, I got error with path "C:\Program ", i.e., has space between Program and Files
        "ps" = "tasklist.exe"
    }
    $systemAliasList = @(
        "cp",
		"cat",
		"mkdir",
		"ls",
		"mv",
		"rm",
		"rmdir",
		"sleep",
		"sort",
		"tee",
		"curl",
		"grep",
		"sed",
        "ps"
    )
    foreach ($aliasKey in $aliasMap.Keys) {
        $aliasValue = $aliasMap[$aliasKey]
        $isShellNoArgs = ""
        if ($aliasKey -eq "now") {
            $isShellNoArgs = "-o"
        }
        alias-manager-cu -n "$aliasKey" -c "$aliasValue" $isShellNoArgs $SCRIPT_FORCE_ARG
    }
    header-log "Disable system alias"
    foreach ($systemAlias in $systemAliasList) {
        disable-system-alias-cu -n "$systemAlias"
    }
}

function Install-MenusEntries {
    create-menu-entry -n Bash -e "$(whichc wt.exe -n)" -a "--title Bash -d `"$(HOME_DIR.exe)`" `"$(whichc bash.cmd -n)`" --login -i" -c "ConsoleOnly;System;" -i "C:\Windows\System32\cmd.exe,0"
    create-menu-entry -n "ChangeDNS" -e "$(whichc change-dns.exe -n)" -i "C:\Windows\System32\cmd.exe,0" -t -A
    create-menu-entry -n "Update System CU" -e "system-upgrade" -a "& pause" -t
    create-menu-entry -n "Cleanup System CU" -e "system-cleanup" -a "& pause" -t
}

switch ($OPERATION_ARG) {
    scripts-apps { Install-ScriptsApps }
    alias { Install-Alias }
    menu-entries { Install-MenusEntries }
    Default { warn-log "No Scripts apps or alias operation to process." }
}
