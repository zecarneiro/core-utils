# Author: José M. C. Noronha
# IMPORTANT: Save this script with UTF-8 with BOM if you have problems with characters

param(
  [int] $style
)

$__PROMPT_GREEN_COLOR__="Green"
$__PROMPT_CYAN_COLOR__="Cyan"
$__PROMPT_DARK_GRAY_COLOR__="DarkGray"
$__PROMPT_RED_COLOR__="Red"

function processBaseStyle($workingDir) {
  $is_admin = "false"
  try {
    $is_admin = (is-admin)
  } catch {
    # Do Nothing
  }
  # Find out if the current user identity is elevated (has admin rights)
  if ("$is_admin" -eq "true") {
    $Host.UI.RawUI.WindowTitle += " [ADMIN]"
  } else {
    $Host.UI.RawUI.WindowTitle = "PowerShell {0}" -f $PSVersionTable.PSVersion.ToString()
  }
}

function style_1($workingDir) {
  # prompt vars
  $arrow = "$([char]0x276F)"
  $unionLineStart="╭─"
  $unionLineEnd="╰─"

  ### Build prompt ###
  processBaseStyle "$workingDir"
  Write-Host "${unionLineStart}[" -ForegroundColor "$__PROMPT_GREEN_COLOR__" -NoNewline
  Write-Host "$workingDir" -ForegroundColor "$__PROMPT_CYAN_COLOR__" -NoNewline
  Write-Host "]" -ForegroundColor "$__PROMPT_GREEN_COLOR__"
  # Prompt
  Write-Host "${unionLineEnd}${arrow}" -ForegroundColor "${__PROMPT_GREEN_COLOR__}" -NoNewline
}

function style_2($workingDir) {
  # prompt vars
  $arrow = "$([char]0x2192)"

  ### Build prompt ###
  processBaseStyle "$workingDir"
  Write-Host "$workingDir" -ForegroundColor "$__PROMPT_CYAN_COLOR__"
  # Prompt
  Write-Host "${arrow}" -ForegroundColor "${__PROMPT_GREEN_COLOR__}" -NoNewline
}

function style_3($workingDir) {
  # prompt vars
  $unionStartChar = "["
  $unionEndChar = "]"
  $username = "$([System.Environment]::UserName)"
  $hostname = "$([System.Environment]::MachineName)"
  $arrow = "$([char]0x276F)"
  $unionLineStart="╭─"
  $unionLineEnd="╰─"

  ### Build prompt ###
  processBaseStyle "$workingDir"
  Write-Host "$unionLineStart" -NoNewline
  # Username
  Write-Host $unionStartChar -ForegroundColor "$__PROMPT_DARK_GRAY_COLOR__" -NoNewline
  Write-Host "$username" -ForegroundColor "$__PROMPT_RED_COLOR__" -NoNewline
  Write-Host $unionEndChar -ForegroundColor "$__PROMPT_DARK_GRAY_COLOR__" -NoNewline
  # Hostname
  Write-Host $unionStartChar -ForegroundColor "$__PROMPT_DARK_GRAY_COLOR__" -NoNewline
  Write-Host "$hostname" -ForegroundColor "$__PROMPT_GREEN_COLOR__" -NoNewline
  Write-Host $unionEndChar -ForegroundColor "$__PROMPT_DARK_GRAY_COLOR__" -NoNewline
  # Working Dir
  Write-Host $unionStartChar -ForegroundColor "$__PROMPT_DARK_GRAY_COLOR__" -NoNewline
  Write-Host "$workingDir" -ForegroundColor "$__PROMPT_CYAN_COLOR__" -NoNewline
  Write-Host $unionEndChar -ForegroundColor "$__PROMPT_DARK_GRAY_COLOR__"
  # Prompt
  Write-Host $unionLineEnd -NoNewline
  Write-Host "${arrow}" -ForegroundColor "${__PROMPT_GREEN_COLOR__}" -NoNewline
}

function style_4($workingDir) {
  # prompt vars
  $arrow = "$([char]0x276F)"

  ### Build prompt ###
  processBaseStyle "$workingDir"
  Write-Host "$workingDir" -ForegroundColor "$__PROMPT_CYAN_COLOR__"
  # Prompt
  Write-Host "${arrow}" -ForegroundColor "${__PROMPT_GREEN_COLOR__}" -NoNewline
}

function base_style($workingDir) {
  ### Build prompt ###
  processBaseStyle "$workingDir"
  # Prompt
  Write-Host "[" -NoNewline
  Write-Host "$workingDir" -ForegroundColor "$__PROMPT_CYAN_COLOR__" -NoNewline
  Write-Host "]"
  Write-Host "$" -NoNewline
}

$workingDir = "$(Get-Location)".replace("${home}",'~')
switch ($style) {
  1 { style_1 "$workingDir" }
  2 { style_2 "$workingDir" }
  3 { style_3 "$workingDir" }
  4 { style_4 "$workingDir" }
  Default { base_style "$workingDir" }
}

