# Author: José M. C. Noronha
# IMPORTANT: Save this script with UTF-8 with BOM if you have problems with characters

# Global Vars
$IS_INIT_PROMPT=$true
$MY_SHELL_PROFILE = $profile.CurrentUserAllHosts
$MY_CUSTOM_SHELL_PROFILE = "$home\.powershell-profile-custom.ps1"
$MY_ALIAS = "$home\.powershell_aliases.ps1"
$CONFIG_DIR = "$home\.config"
$OTHER_APPS_DIR = "$home\.otherapps"
$USER_BIN_DIR = "$home\.local\bin"

# BASH-LIKE TAB COMPLETION IN POWERSHELL
Set-PSReadlineKeyHandler -Key Tab -Function Complete
Set-PSReadLineKeyHandler -Key UpArrow -Function HistorySearchBackward
Set-PSReadLineKeyHandler -Key DownArrow -Function HistorySearchForward

$OutputEncoding = [Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8
function isadmin {
  $currentUser = New-Object Security.Principal.WindowsPrincipal $([Security.Principal.WindowsIdentity]::GetCurrent())
  return ($currentUser.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator))
}

# Find out if the current user identity is elevated (has admin rights)
$Host.UI.RawUI.WindowTitle = "PowerShell {0}" -f $PSVersionTable.PSVersion.ToString()
try {
  if ((isadmin)) {
    $Host.UI.RawUI.WindowTitle += " [ADMIN]"
  }
} catch {
  errorlog "An error occurred, when tray to check if is admin and set [ADMIN] on title of windows terminal"
}

function prompt {
  $hasNoError = $?
  if ((isadmin)) {
    "[" + (Get-Location) + "] # "
  }
  else {
    $promptChar = "$([char]0x276F)"
    $promptCharColor = "Green"
    $actualLocation = "$(Get-Location)".replace("${home}",'~')
    if (-not $hasNoError) {
      # Last command failed
      $promptCharColor = "Red"
    }
    if ($global:IS_INIT_PROMPT) {
      $global:IS_INIT_PROMPT=$false
    } else {
      $actualLocation = "`n$actualLocation"
    }
    Write-Host "$actualLocation" -ForegroundColor "Cyan"
    Write-Host "${promptChar}" -ForegroundColor "${promptCharColor}" -NoNewline
    " "
  }
}
