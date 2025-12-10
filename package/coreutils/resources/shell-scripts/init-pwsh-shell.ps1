# Author: José M. C. Noronha
# IMPORTANT: Save this script with UTF-8 with BOM if you have problems with characters

# Global Vars
$__IS_INIT_PROMPT__=$true
$__COREUTILS_SCRIPT_DIR__="$PSScriptRoot"
$__COREUTILS_LIBS_SCRIPT_DIR__="${__COREUTILS_SCRIPT_DIR__}\libs"

# If is on Windows Server OS, change security protocol
if ((OS_NAME) -contains "Server") {
  [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
}

$OutputEncoding = [Console]::OutputEncoding = New-Object System.Text.UTF8Encoding $False

# IMPORT LIBS FILES
$__LIBS_TO_NOT_IMPORT__=@("style-prompt.ps1", "trash.ps1", "pgrep.ps1", "pkill.ps1", "has-internet.ps1", "download.ps1", "run-line-as-command.ps1")
Get-ChildItem -Path "$(Resolve-Path "${__COREUTILS_LIBS_SCRIPT_DIR__}")" -Filter *.ps1 -File | ForEach-Object {
  $lib_to_import = $_.FullName
  $can_import = $true 
  foreach ($lib in $__LIBS_TO_NOT_IMPORT__) {
    if ("$lib_to_import".Contains("${lib}")) {
      $can_import = $false
      break
    }
  }
  if ($can_import) {
    . "$lib_to_import"
  }
}

# SHELL FUNCTIONS
try {
    $resolved_path = (Resolve-Path("$home\.local\coreutils\scripts\pwsh-shell") -ErrorAction SilentlyContinue)
    $cmd_path = (Resolve-Path("$resolved_path\cmd") -ErrorAction SilentlyContinue)
    if (Test-Path "$cmd_path") {
        $env:Path += ";$cmd_path"
    } else {
        $env:PATH = "${resolved_path}:${env:PATH}"
    }
} catch {
  # Do Nothing
}

# IMPORT ALIAS
try {
  . "$(Resolve-Path("$home\.local\coreutils\alias\powershell-alias.ps1") -ErrorAction SilentlyContinue)"
} catch {
  # Do Nothing
}
function .. { Set-Location .. }

# DEPENDENCIES
# ENABLE READLINE
try {
  Set-PSReadlineKeyHandler -Key Tab -Function Complete
  Set-PSReadLineKeyHandler -Key UpArrow -Function HistorySearchBackward
  Set-PSReadLineKeyHandler -Key DownArrow -Function HistorySearchForward
} catch {
  # Do Nothing
}
# ENABLE FZF
try {
  Set-PsFzfOption -EnableFd:$true
  Set-PsFzfOption -PSReadlineChordProvider 'Ctrl+t' -PSReadlineChordReverseHistory 'Ctrl+r'
} catch {
  # Do Nothing
}

function prompt {
  $style_script = "$($global:__COREUTILS_LIBS_SCRIPT_DIR__)\style-prompt.ps1"
  $style_script = "$(Resolve-Path($style_script) -ErrorAction SilentlyContinue)"
  . "${style_script}" -style (prompt-style -s) -is_init $global:__IS_INIT_PROMPT__
  if ($global:__IS_INIT_PROMPT__) {
    $global:__IS_INIT_PROMPT__=$false
  }
}
