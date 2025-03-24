#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2154

# Global Vars
declare IS_INIT_PROMPT=true

function isadmin {
    if [ "$(id -u)" -eq 0 ]; then
        echo true
    else
        echo false
    fi
}

build_prompt() {
    # prompt vars
    unionStartChar="${DarkGrayColor}[${BoldColor}${ResetColor}"
    unionEndChar="${DarkGrayColor}]${BoldColor}${ResetColor}"
    username="${unionStartChar}${RedColor}\u${ResetColor}${unionEndChar}"
    hostname="${unionStartChar}${GreenColor}\h${ResetColor}${unionEndChar}"
    workingDir="${unionStartChar}${CyanColor}\w${ResetColor}${unionEndChar}"
    arrow="${GreenColor}❯${BoldColor}${ResetColor}"
    unionLineStart="╭─"
    unionLineEnd="╰─"
    windowsTitle="\[\e]0;Bash \v\a\]"

    ### Build prompt ###
    promptBuilder="${unionLineStart}${username}${hostname}${workingDir}\n${unionLineEnd}${arrow} "
    if [ "$IS_INIT_PROMPT" == "false" ]; then
        promptBuilder="\n${promptBuilder}"
    else
        IS_INIT_PROMPT=false
    fi
    PS1="${debian_chroot:+($debian_chroot)}${windowsTitle}${promptBuilder}"  # begin prompt
}
PROMPT_COMMAND=build_prompt
