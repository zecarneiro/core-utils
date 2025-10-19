# Author: José M. C. Noronha
# IMPORTANT: Save this script with UTF-8 with BOM if you have problems with characters

# Global Vars
$IS_INIT_PROMPT=$true
$IS_ADMIN_ROLE=$false

# If is on Windows Server OS, change security protocol
if ($MY_OS -contains "Server") {
  [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
}

# ENABLE READLINE
try {
  Set-PSReadlineKeyHandler -Key Tab -Function Complete
  Set-PSReadLineKeyHandler -Key UpArrow -Function HistorySearchBackward
  Set-PSReadLineKeyHandler -Key DownArrow -Function HistorySearchForward
}
catch {
  # Do Nothing
}

# ENABLE FZF
try {
  Set-PsFzfOption -EnableFd:$true
  Set-PsFzfOption -PSReadlineChordProvider 'Ctrl+t' -PSReadlineChordReverseHistory 'Ctrl+r'
}
catch {
  # Do Nothing
}

$OutputEncoding = [Console]::OutputEncoding = New-Object System.Text.UTF8Encoding $False
function isadmin {
  $currentUser = New-Object Security.Principal.WindowsPrincipal $([Security.Principal.WindowsIdentity]::GetCurrent())
  return ($currentUser.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator))
}

try {
  $IS_ADMIN_ROLE = (isadmin)
} catch {
  $IS_ADMIN_ROLE = $false
}
if ($IS_ADMIN_ROLE) {
  $Host.UI.RawUI.WindowTitle += " [ADMIN]"
}

# Find out if the current user identity is elevated (has admin rights)
$Host.UI.RawUI.WindowTitle = "PowerShell {0}" -f $PSVersionTable.PSVersion.ToString()
function prompt {
  if ($IS_ADMIN_ROLE) {
    "[" + (Get-Location) + "] # "
  }
  else {
    # prompt vars
    $workingDir = "$(Get-Location)".replace("${home}",'~')
    $arrow = "$([char]0x276F)"
    $unionLineStart="╭─"
    $unionLineEnd="╰─"
    
    ### Build prompt ###
    if (!$global:IS_INIT_PROMPT) {
      Write-Host ""
    } else {
      $global:IS_INIT_PROMPT=$false
    }
    # Working Dir
    Write-Host "${unionLineStart}[" -ForegroundColor "$GreenColor" -NoNewline
    Write-Host "$workingDir" -ForegroundColor "$CyanColor" -NoNewline
    Write-Host "]" -ForegroundColor "$GreenColor"
    # Prompt
    Write-Host "${unionLineEnd}${arrow}" -ForegroundColor "${GreenColor}" -NoNewline
    " "
  }
}
