#!/usr/bin/env bash
# Author: José M. C. Noronha

declare MY_SHELL_PROFILE="$HOME/.bashrc"
declare MY_CUSTOM_SHELL_PROFILE="$HOME/.bash-profile-custom.sh"
declare MY_ALIAS="$HOME/.bash_aliases"
declare CONFIG_DIR="$HOME/.config"
declare OTHER_APPS_DIR="$HOME/.otherapps"
declare USER_BIN_DIR="$HOME/.local/bin"

# COLORS
declare ResetColor='\033[0m'           # Text Reset
declare BoldColor='\033[0;1m'
declare BackgroundColor="\033[0;101m"
# Regular Colors
declare BlackColor='\033[0;30m'        # Black
declare RedColor='\033[0;31m'          # Red
declare GreenColor='\033[0;32m'        # Green
declare YellowColor='\033[0;33m'       # Yellow
declare BlueColor='\033[0;34m'         # Blue
declare PurpleColor='\033[0;35m'       # Purple
declare CyanColor='\033[0;36m'         # Cyan
declare WhiteColor='\033[0;37m'        # White
declare DarkGrayColor='\033[1;30m'     # Dark Gray
declare GrayColor='\033[90m'           # Gray
