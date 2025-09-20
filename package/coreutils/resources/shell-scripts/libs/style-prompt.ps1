param(
  [int] $style,
  [bool] $is_init
)

$__PROMPT_GREEN_COLOR__="Green"
$__PROMPT_CYAN_COLOR__="Cyan"
$__PROMPT_DARK_GRAY_COLOR__="DarkGray"
$__PROMPT_RED_COLOR__="Red"

function processBaseStyle($workingDir) {
  # Find out if the current user identity is elevated (has admin rights)
  if ((is-admin) -eq "true") {
    $Host.UI.RawUI.WindowTitle += " [ADMIN]"
    "[" + $workingDir + "] # "
  } else {
    $Host.UI.RawUI.WindowTitle = "PowerShell {0}" -f $PSVersionTable.PSVersion.ToString()
  }
  if (!$is_init) {
    Write-Host ""
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
  " "
}

function style_2($workingDir) {
  # prompt vars
  $arrow = "$([char]0x2192)"

  ### Build prompt ###
  processBaseStyle "$workingDir"
  Write-Host "$workingDir" -ForegroundColor "$__PROMPT_CYAN_COLOR__"
  # Prompt
  Write-Host "${arrow}" -ForegroundColor "${__PROMPT_GREEN_COLOR__}" -NoNewline
  " "
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
  " "
}

function style_4($workingDir) {
  # prompt vars
  $arrow = "$([char]0x276F)"

  ### Build prompt ###
  processBaseStyle "$workingDir"
  Write-Host "$workingDir" -ForegroundColor "$__PROMPT_CYAN_COLOR__"
  # Prompt
  Write-Host "${arrow}" -ForegroundColor "${__PROMPT_GREEN_COLOR__}" -NoNewline
  " "
}

$workingDir = "$(Get-Location)".replace("${home}",'~')
switch ($style) {
  1 { style_1 "$workingDir" }
  2 { style_2 "$workingDir" }
  3 { style_3 "$workingDir" }
  4 { style_4 "$workingDir" }
  Default { Write-Host "$(Split-Path -Leaf $MyInvocation.MyCommand.Path) style(Accept 1..4) is_init(Accept true|false)" }
}

