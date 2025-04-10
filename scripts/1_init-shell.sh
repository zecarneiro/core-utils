#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2154

# Global Vars
declare IS_INIT_PROMPT=true
declare FZF_KEY_BINDINGS_EXAMPLE_FILE="/usr/share/doc/fzf/examples/key-bindings.bash"
declare FZF_KEY_BINDINGS_COMPLETION_FILE="/usr/share/bash-completion/completions/fzf-key-bindings"

# ENABLE FZF
if [ -f "$FZF_KEY_BINDINGS_EXAMPLE_FILE" ]; then
    # shellcheck source=/dev/null
    source "$FZF_KEY_BINDINGS_EXAMPLE_FILE"
fi
if [ -f "$FZF_KEY_BINDINGS_COMPLETION_FILE" ]; then
    # shellcheck source=/dev/null
    source "$FZF_KEY_BINDINGS_COMPLETION_FILE"
fi
# Do if on Windows SO
if [[ "$OSTYPE" == "cygwin" ]]||[[ "$OSTYPE" == "msys" ]]||[[ "$OSTYPE" == "win32" ]]; then
    export LANG=C.UTF-8
fi

function isadmin {
    if [ "$(id -u)" -eq 0 ]; then
        echo true
    else
        echo false
    fi
}

build_prompt() {
    # colors vars
    local boldColor="\[${BoldColor}\]"
    local greenColor="\[${GreenColor}\]"
    local resetColor="\[${ResetColor}\]"
    local cyanColor="\[${CyanColor}\]"

    # prompt vars
    local workingDir="${cyanColor}\w${resetColor}"
    local arrow="${greenColor}❯${boldColor}${resetColor}"
    local windowsTitle="\[\e]0;Bash \v\a\]"

    ### Build prompt ###
    local promptBuilder="${workingDir}\n${arrow} "
    if [ "$IS_INIT_PROMPT" == "false" ]; then
        promptBuilder="\n${promptBuilder}"
    else
        IS_INIT_PROMPT=false
    fi
    PS1="${debian_chroot:+($debian_chroot)}${windowsTitle}${promptBuilder}"  # begin prompt
}
PROMPT_COMMAND=build_prompt
