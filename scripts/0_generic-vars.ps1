# Author: José M. C. Noronha

# --------------------------------- BASE VARS -------------------------------- #
$MY_SHELL_PROFILE = $profile.CurrentUserAllHosts
$MY_CUSTOM_SHELL_PROFILE = "$home\.powershell-profile-custom.ps1"
$MY_ALIAS = "$home\.powershell_aliases.ps1"

# ----------------------------------- DIRS ----------------------------------- #
$CONFIG_DIR = "$home\.config"
$OTHER_APPS_DIR = "$home\.otherapps"
$USER_BIN_DIR = "$home\.local\bin"
$USER_STARTUP_DIR = "${home}\Start Menu\Programs\Startup"
$TEMP_DIR="$([System.IO.Path]::GetTempPath())".TrimEnd("\")

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
