# Author: José M. C. Noronha
# IMPORTANT: Save this script with UTF-8 with BOM if you have problems with characters

# Global Vars
$__COREUTILS_SCRIPT_DIR__="$PSScriptRoot"
$__COREUTILS_LIBS_SCRIPT_DIR__="${__COREUTILS_SCRIPT_DIR__}\..\libs"

# If is on Windows Server OS, change security protocol
try {
  if ($PSVersionTable.Platform -ne "Unix") {
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]'Tls,Tls11,Tls12'
  }
} catch {
  # Do Nothing
}

$OutputEncoding = [Console]::OutputEncoding = New-Object System.Text.UTF8Encoding $False

# Set Others Alias
function .. { Set-Location .. }

# ENABLE READLINE
try {
  Set-PSReadlineKeyHandler -Key Tab -Function Complete
  Set-PSReadLineKeyHandler -Key UpArrow -Function HistorySearchBackward
  Set-PSReadLineKeyHandler -Key DownArrow -Function HistorySearchForward
} catch {
  # Failed to ENABLE READLINE
}

# ENABLE FZF
try {
  Set-PsFzfOption -EnableFd:$true
  Set-PsFzfOption -PSReadlineChordProvider 'Ctrl+t' -PSReadlineChordReverseHistory 'Ctrl+r'
} catch {
  # Failed to ENABLE FZF
}

if (Test-Path "$home\.local\share\coreutils\system-aliases\system-alias-pwsh.ps1") { . "$home\.local\share\coreutils\system-aliases\system-alias-pwsh.ps1" }

function prompt {
  $style_to_use = -1
  if ([System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform([System.Runtime.InteropServices.OSPlatform]::Windows)) {
    if (Get-Command prompt-style -ErrorAction SilentlyContinue) {
      $style_to_use = (prompt-style -s)
    }    
  }
  $style_script = "$($global:__COREUTILS_LIBS_SCRIPT_DIR__)\style-prompt.ps1"
  $style_script = "$(Resolve-Path($style_script) -ErrorAction SilentlyContinue)"
  . "${style_script}" -style $style_to_use
  return " "
}

