#!/usr/bin/env bash
# Author: José M. C. Noronha

declare style="$1"
declare is_init="$2"
declare __PROMPT_RESET_COLOR__='\033[0m'           # Text Reset
declare __PROMPT_BOLD_COLOR__='\033[0;1m'
declare __PROMPT_GREEN_COLOR__='\033[0;32m'
declare __PROMPT_CYAN_COLOR__='\033[0;36m'
declare __PROMPT_DARK_GRAY_COLOR__='\033[1;30m'
declare __PROMPT_RED_COLOR__='\033[0;31m'

style_1() {
    # colors vars
    local resetColor="\[${__PROMPT_RESET_COLOR__}\]"
    local boldColor="\[${__PROMPT_BOLD_COLOR__}\]"
    local greenColor="\[${__PROMPT_GREEN_COLOR__}\]"
    local cyanColor="\[${__PROMPT_CYAN_COLOR__}\]"
    # prompt vars
    local workingDir="${cyanColor}\w${resetColor}"
    local arrow="${greenColor}❯${boldColor}${resetColor}"
    local unionLineStart="${greenColor}╭─[${resetColor}"
    local unionLineEnd="${greenColor}╰─${resetColor}"
    ### Build prompt ###
    echo -e "${unionLineStart}${workingDir}${greenColor}]${resetColor}\n${unionLineEnd}${arrow} "
}

style_2() {
    # colors vars
    local resetColor="\[${__PROMPT_RESET_COLOR__}\]"
    local boldColor="\[${__PROMPT_BOLD_COLOR__}\]"
    local greenColor="\[${__PROMPT_GREEN_COLOR__}\]"
    local cyanColor="\[${__PROMPT_CYAN_COLOR__}\]"
    # prompt vars
    local workingDir="${cyanColor}\w${resetColor}"
    local arrow="${greenColor}→${boldColor}${resetColor}"
    ### Build prompt ###
    echo -e "${workingDir}\n${arrow} "
}

style_3() {
    # colors vars
    local darkGrayColor="\[${__PROMPT_DARK_GRAY_COLOR__}\]"
    local boldColor="\[${__PROMPT_BOLD_COLOR__}\]"
    local redColor="\[${__PROMPT_RED_COLOR__}\]"
    local greenColor="\[${__PROMPT_GREEN_COLOR__}\]"
    local resetColor="\[${__PROMPT_RESET_COLOR__}\]"
    local cyanColor="\[${__PROMPT_CYAN_COLOR__}\]"
    # prompt vars
    local unionStartChar="${darkGrayColor}[${boldColor}${resetColor}"
    local unionEndChar="${darkGrayColor}]${boldColor}${resetColor}"
    local username="${unionStartChar}${redColor}\u${resetColor}${unionEndChar}"
    local hostname="${unionStartChar}${greenColor}\h${resetColor}${unionEndChar}"
    local workingDir="${unionStartChar}${cyanColor}\w${resetColor}${unionEndChar}"
    local arrow="${greenColor}❯${boldColor}${resetColor}"
    local unionLineStart="╭─"
    local unionLineEnd="╰─"
    ### Build prompt ###
    echo -e "${unionLineStart}${username}${hostname}${workingDir}\n${unionLineEnd}${arrow} "
}

style_4() {
    # colors vars
    local boldColor="\[${__PROMPT_BOLD_COLOR__}\]"
    local greenColor="\[${__PROMPT_GREEN_COLOR__}\]"
    local resetColor="\[${__PROMPT_RESET_COLOR__}\]"
    local cyanColor="\[${__PROMPT_CYAN_COLOR__}\]"
    # prompt vars
    local workingDir="${cyanColor}\w${resetColor}"
    local arrow="${greenColor}❯${boldColor}${resetColor}"
    ### Build prompt ###
    echo -e "${workingDir}\n${arrow} "
}

declare windows_title="\[\e]0;Bash \v\a\]"
declare prompt_builder=""
case "${style}" in
    1) prompt_builder="$(style_1)" ;;
    2) prompt_builder="$(style_2)" ;;
    3) prompt_builder="$(style_3)" ;;
    4) prompt_builder="$(style_4)" ;;
    *) echo "$0 style(Accept 1..4) is_init(Accept true|false)" ;;
esac
if [[ -n "${prompt_builder}" ]]; then
    if [ "$is_init" == "false" ]; then
        prompt_builder="\n${prompt_builder}"
    fi
    PS1="${debian_chroot:+($debian_chroot)}${windows_title}${prompt_builder}"  # begin prompt
fi
