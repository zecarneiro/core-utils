# Author: José M. C. Noronha

# --------------------------------- BASE VARS -------------------------------- #
try {
    $__osInfo = Get-CimInstance -ClassName Win32_OperatingSystem
    $MY_OS=$__osInfo.Caption
}
catch {
    $MY_OS=""
}
$MY_SHELL_PROFILE = $profile.CurrentUserAllHosts
$MY_CUSTOM_SHELL_PROFILE = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath("$home\.powershell-profile-custom.ps1")
$MY_ALIAS = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath("$home\.powershell_aliases.ps1")

# ----------------------------------- DIRS ----------------------------------- #
$CONFIG_DIR = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath("$home\.config")
$OTHER_APPS_DIR = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath("$home\.otherapps")
$USER_BIN_DIR = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath("$home\.local\bin")
$USER_STARTUP_DIR = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath("${home}\Start Menu\Programs\Startup")
$USER_TEMP_DIR = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath("$([System.IO.Path]::GetTempPath())".TrimEnd("\").TrimEnd("/"))
$TEMP_DIR = "$USER_TEMP_DIR"

# ---------------------------------- COLORS ---------------------------------- #
# User Write-Host to print colored text
# BoldColor -> Use -ForegroundColor
# BackgroundColor -> Use -BackgroundColor
# Regular Colors
$BlackColor="Black"
$RedColor="Red"
$GreenColor="Green"
$YellowColor="Yellow"
$BlueColor="Blue"
$PurpleColor="Purple"
$CyanColor="Cyan"
$WhiteColor="White"
$DarkGrayColor="DarkGray"
$GrayColor="Gray"
