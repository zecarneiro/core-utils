#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2154
# shellcheck disable=SC1090
# shellcheck disable=SC2155

# Global Vars
declare __IS_INIT_PROMPT__=true
declare __COREUTILS_SCRIPT_DIR__=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
declare __COREUTILS_LIBS_SCRIPT_DIR__="${__COREUTILS_SCRIPT_DIR__}/libs"

# Do if on Windows SO
if [[ "$OSTYPE" == "cygwin" ]]||[[ "$OSTYPE" == "msys" ]]||[[ "$OSTYPE" == "win32" ]]; then
    export LANG=C.UTF-8
fi

# IMPORT LIBS FILES
declare -a __LIBS_TO_NOT_IMPORT__=("has-internet.sh" "run-line-as-command.sh" "style-prompt.sh")
for lib_to_import in "${__COREUTILS_LIBS_SCRIPT_DIR__}"/*.sh; do
  can_import=true
  for lib in "${__LIBS_TO_NOT_IMPORT__[@]}"; do
    if echo "$lib_to_import" | grep -q "$lib"; then
      can_import=false
      break
    fi
  done
  if [[ "${can_import}" == true ]]; then
    . "$lib_to_import"
  fi
done

# SHELL SCRIPTS
export PATH="$HOME/.local/coreutils/scripts/linux-shell:$PATH"

# IMPORT ALIAS
if [ -f "$HOME/.local/coreutils/alias/bash-alias" ]; then
  source "$HOME/.local/coreutils/alias/bash-alias"
fi
alias ..='cd ..'

# DEPENDENCY AREA
eval "$(fzf --bash)" # ENABLE FZF

build_prompt() {
  local style_script="${__COREUTILS_LIBS_SCRIPT_DIR__}/style-prompt.sh"
  . "${style_script}" "$(prompt-style -s)" "$__IS_INIT_PROMPT__"
  if [[ "$__IS_INIT_PROMPT__" == "true" ]]; then
    __IS_INIT_PROMPT__=false
  fi
}
PROMPT_COMMAND=build_prompt
