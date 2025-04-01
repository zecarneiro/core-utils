# Author: José M. C. Noronha
# IMPORTANT: Save this script with UTF-8 with BOM if you have problems with characters

# Global Vars
$IS_INIT_PROMPT=$true
$IS_ADMIN_ROLE=$false

# ENABLE FZF AND READLINE
try {
  Set-PsFzfOption -EnableFd:$true
  Set-PSReadlineKeyHandler -Key Tab -Function Complete
  Set-PSReadLineKeyHandler -Key UpArrow -Function HistorySearchBackward
  Set-PSReadLineKeyHandler -Key DownArrow -Function HistorySearchForward
  Set-PsFzfOption -PSReadlineChordProvider 'Ctrl+t' -PSReadlineChordReverseHistory 'Ctrl+r'
}
catch {
  # Do Nothing
}

$OutputEncoding = [Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8
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
    $unionStartChar = "["
    $unionEndChar = "]"
    $username = "$([System.Environment]::UserName)"
    $hostname = "$([System.Environment]::MachineName)"
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
    Write-Host "$unionLineStart" -NoNewline
    # Username
    Write-Host $unionStartChar -ForegroundColor "$DarkGrayColor" -NoNewline
    Write-Host "$username" -ForegroundColor "$RedColor" -NoNewline
    Write-Host $unionEndChar -ForegroundColor "$DarkGrayColor" -NoNewline
    # Hostname
    Write-Host $unionStartChar -ForegroundColor "$DarkGrayColor" -NoNewline
    Write-Host "$hostname" -ForegroundColor "$GreenColor" -NoNewline
    Write-Host $unionEndChar -ForegroundColor "$DarkGrayColor" -NoNewline
    # Working Dir
    Write-Host $unionStartChar -ForegroundColor "$DarkGrayColor" -NoNewline
    Write-Host "$workingDir" -ForegroundColor "$CyanColor" -NoNewline
    Write-Host $unionEndChar -ForegroundColor "$DarkGrayColor"
    # Prompt
    Write-Host $unionLineEnd -NoNewline
    Write-Host "${arrow}" -ForegroundColor "${GreenColor}" -NoNewline
    " "
  }
}
