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
    
    ### Build prompt ###
    if (!$global:IS_INIT_PROMPT) {
      Write-Host ""
    } else {
      $global:IS_INIT_PROMPT=$false
    }
    # Working Dir
    Write-Host "$workingDir" -ForegroundColor "$CyanColor"
    # Prompt
    Write-Host "${arrow}" -ForegroundColor "${GreenColor}" -NoNewline
    " "
  }
}
